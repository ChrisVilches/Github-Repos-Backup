package github

import (
	"encoding/json"
	"fmt"
	"github-backup-repos/models"
	"net/http"
	"strconv"
)

const perPage = 30

func GetAllRepos(token string) ([]models.Repo, error) {
	allRepos := []models.Repo{}

	url := "https://api.github.com/user/repos"

	for page := 1; ; page++ {
		req, err := http.NewRequest("GET", url, nil)

		if err != nil {
			return nil, err
		}

		q := req.URL.Query()
		q.Add("per_page", strconv.Itoa(perPage))
		q.Add("page", strconv.Itoa(page))
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

		fmt.Println("Fetched page", page, "with", len(repos), "repos")

		allRepos = append(allRepos, repos...)

		if len(repos) < perPage {
			break
		}
	}

	return allRepos, nil
}
