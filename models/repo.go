package models

import (
	"time"
)

type Repo struct {
	Owner       User      `json:"owner"`
	Description string    `json:"description"`
	GitURL      string    `json:"git_url"`
	HTMLURL     string    `json:"html_url"`
	SSHURL      string    `json:"ssh_url"`
	CloneURL    string    `json:"clone_url"`
	Name        string    `json:"name"`
	UpdatedAt   time.Time `json:"updated_at"`
	Private     bool      `json:"private"`
}

func (r Repo) GetID() string {
	return r.Name
}
