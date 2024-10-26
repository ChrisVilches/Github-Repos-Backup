package github

import (
	"fmt"
	"os/exec"
	"path"
)

func CloneRepo(username, repoName, token, destDir string) {
	gitURL := fmt.Sprintf("https://%s:%s@github.com/%s/%s.git", username, token, username, repoName)

	dest := path.Join(destDir, repoName)
	cmd := exec.Command("git", "clone", gitURL, dest)

	err := cmd.Run()

	if err != nil {
		fmt.Printf("Error running command: %v\n", err)
	}
}
