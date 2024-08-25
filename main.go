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

func zipProcess(input, output chan models.Repo, destDir string) {
	for repo := range input {
		src := path.Join(destDir, repo.Name)
		dest := path.Join(destDir, fmt.Sprintf("%s.zip", repo.Name))
		ZipRepo(src, dest, true)
		output <- repo
	}
}

func gitCloneProcess(input, output chan models.Repo, token, cloneDestDir string) {
	for repo := range input {
		fmt.Println("Cloning...", repo.Name)
		github.CloneRepo(repo.Owner.Login, repo.Name, token, cloneDestDir)
		output <- repo
	}
}

func backupRepos(username, token string, numWorkers int, destDir string) {
	fmt.Println(time.Now())
	allRepos := github.GetAllRepos(username, token)
	fmt.Println("Total:", len(allRepos))

	currentSaved, err := util.ReadJSON[models.Repo](path.Join(destDir, "updated-at.json"))

	if err != nil {
		fmt.Println("Error reading current saved repos (doing back-up from scratch):", err)
		currentSaved = []models.Repo{}
	}

	allRepos = filterOwnerRepos(allRepos, username)
	allRepos = filterUpToDateRepos(allRepos, currentSaved)

	fmt.Println("After filtering:", len(allRepos))

	completionCh := make(chan models.Repo)
	gitCloneCh := make(chan models.Repo, len(allRepos))
	zipCh := make(chan models.Repo, len(allRepos))

	for i := 0; i < numWorkers; i++ {
		go zipProcess(zipCh, completionCh, destDir)
		go gitCloneProcess(gitCloneCh, zipCh, token, destDir)
	}

	for _, repo := range allRepos {
		gitCloneCh <- repo
	}

	// This works for synchronization, so it's not necessary to use a WaitGroup.
	for i := 0; i < len(allRepos); i++ {
		repo := <-completionCh
		fmt.Printf("(%d/%d) Completed %s\n", i+1, len(allRepos), repo.Name)
	}

	util.WriteJSON(path.Join(destDir, "updated-at.json"), util.PatchList(currentSaved, allRepos))
	fmt.Println(time.Now())
	fmt.Println("Backup complete:", destDir)
}

// TODO: I think the output folder must be parameterizable, because I have no idea
// how I'm going to distribute this software... is it going to be its own github repo? or will
// it be part of the configs repo?? I think it can be its own repo because it contains nothing
// sensitive
// TODO: should be able to work without token... that way I just clone the public repos.
func main() {
	username := flag.String("username", "", "GitHub username")
	token := flag.String("token", "", "GitHub API token")
	destDir := flag.String("dest-dir", "./repos", "Destination directory for cloned repos")

	flag.Parse()

	if *destDir == "" {
		fmt.Println("Destination directory is required.")
		flag.Usage()
		os.Exit(1)
	}

	if *username == "" || *token == "" {
		fmt.Println("Both username and token are required.")
		flag.Usage()
		os.Exit(1)
	}

	numWorkers := 10

	finalPath, err := filepath.Abs(*destDir)

	if err != nil {
		fmt.Println("Error getting absolute path:", err)
		os.Exit(1)
	}

	backupRepos(*username, *token, numWorkers, finalPath)
}
