package main

import (
	"flag"
	"fmt"
	"github-backup-repos/github"
	"github-backup-repos/models"
	"github-backup-repos/util"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"
)

func filterUpToDateRepos(repos []models.Repo, currentSaved []models.Repo) []models.Repo {
	var result []models.Repo

	savedMap := make(map[string]models.Repo)

	for _, repo := range currentSaved {
		savedMap[repo.Name] = repo
	}

	for _, repo := range repos {
		savedRepo, ok := savedMap[repo.Name]

		if (ok && savedRepo.UpdatedAt != repo.UpdatedAt) || !ok {
			fmt.Println("adding repo", repo.Name, repo.UpdatedAt, savedRepo.UpdatedAt)
			result = append(result, repo)
		}
	}

	return result
}

func filterOwnerRepos(repos []models.Repo, username string) []models.Repo {
	var result []models.Repo

	for _, repo := range repos {
		if repo.Owner.Login == username {
			result = append(result, repo)
		}
	}

	return result
}

func zipWork(repo models.Repo, destDir string) {
	src := path.Join(destDir, repo.Name)
	dest := path.Join(destDir, fmt.Sprintf("%s.zip", repo.Name))
	err := zipRepo(src, dest, true)
	if err != nil {
		panic(err)
	}
}

func zipProcess(jobs <-chan models.Repo, numWorkers int, destDir string) <-chan models.Repo {
	zipped := make(chan models.Repo)
	var wg sync.WaitGroup

	for range numWorkers {
		wg.Add(1)
		go func() {
			for repo := range jobs {
				zipWork(repo, destDir)
				zipped <- repo
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(zipped)
	}()

	return zipped
}

func gitCloneProcess(
	jobs <-chan models.Repo,
	numWorkers int,
	token,
	cloneDestDir string,
) <-chan models.Repo {
	cloned := make(chan models.Repo)
	var wg sync.WaitGroup

	for range numWorkers {
		wg.Add(1)
		go func() {
			for repo := range jobs {
				fmt.Println("Cloning...", repo.Name)
				github.CloneRepo(repo.Owner.Login, repo.Name, token, cloneDestDir)
				cloned <- repo
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(cloned)
	}()

	return cloned
}

func backupRepos(username, token string, numWorkers int, destDir string) error {
	fmt.Println(time.Now())
	allRepos, err := github.GetAllRepos(token)
	if err != nil {
		return err
	}

	fmt.Println("Total:", len(allRepos))

	currentSaved, err := util.ReadJSON[models.Repo](path.Join(destDir, "updated-at.json"))

	if err != nil {
		fmt.Println("Error reading current saved repos (doing back-up from scratch):", err)
		currentSaved = []models.Repo{}
	}

	allRepos = filterOwnerRepos(allRepos, username)
	allRepos = filterUpToDateRepos(allRepos, currentSaved)

	fmt.Println("After filtering:", len(allRepos))

	jobs := util.ListToReadonlyChannel(allRepos, 0)

	completionCh := zipProcess(
		gitCloneProcess(jobs, numWorkers, token, destDir),
		numWorkers,
		destDir,
	)

	idx := 1
	for repo := range completionCh {
		fmt.Printf("(%d/%d) Completed %s\n", idx, len(allRepos), repo.Name)
		idx++
	}

	util.WriteJSON(path.Join(destDir, "updated-at.json"), util.PatchList(currentSaved, allRepos))
	fmt.Println(time.Now())
	fmt.Println("Backup complete:", destDir)

	return nil
}

const (
	errorStatusCode = 1
	numWorkers      = 10
)

// TODO: should be able to work without token... that way I just clone the public repos.
func main() {
	username := flag.String("username", "", "GitHub username")
	token := flag.String("token", "", "GitHub API token")
	destDir := flag.String("dest-dir", "./repos", "Destination directory for cloned repos")

	flag.Parse()

	if *destDir == "" {
		fmt.Println("Destination directory is required.")
		flag.Usage()
		os.Exit(errorStatusCode)
	}

	if *username == "" || *token == "" {
		fmt.Println("Both username and token are required.")
		flag.Usage()
		os.Exit(errorStatusCode)
	}

	finalPath, err := filepath.Abs(*destDir)

	if err != nil {
		fmt.Println("Error getting absolute path:", err)
		os.Exit(errorStatusCode)
	}

	err = backupRepos(*username, *token, numWorkers, finalPath)

	if err != nil {
		fmt.Println(err)
		os.Exit(errorStatusCode)
	}
}
