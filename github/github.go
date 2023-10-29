package github

import "github.com/google/go-github/v56/github"

func ProvideGithub() *github.Client {
	return github.NewClient(nil)
}

var Options = ProvideGithub
