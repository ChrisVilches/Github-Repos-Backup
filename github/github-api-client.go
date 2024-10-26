package github

import (
	"encoding/json"
	"fmt"
	"github-backup-repos/models"
	"net/http"
	"strconv"
)

const (
	perPage = 30
	apiURL  = "https://api.github.com/user/repos"
)

func getReposPage(pageIdx int, token string) ([]models.Repo, error) {
	req, err := http.NewRequest("GET", apiURL, nil)

	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("per_page", strconv.Itoa(perPage))
	q.Add("page", strconv.Itoa(pageIdx))
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Authorization", fmt.Sprintf("token %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var repos []models.Repo
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		fmt.Printf("Error decoding response: %s\n", err)
		return nil, err
	}

	fmt.Println("Fetched page", pageIdx, "with", len(repos), "repos")

	return repos, nil
}

func GetAllRepos(token string) ([]models.Repo, error) {
	allRepos := []models.Repo{}

	for page := 1; ; page++ {
		repos, err := getReposPage(page, token)

		if err != nil {
			return nil, err
		}

		allRepos = append(allRepos, repos...)

		if len(repos) < perPage {
			break
		}
	}

	return allRepos, nil
}
