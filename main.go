package main

import (
	"context"
	"log"

	"github.com/kelseyhightower/envconfig"
)

type env struct {
	GithubToken     string `envconfig:"GITHUB_TOKEN"`
	ReleaseFilePath string `envconfig:"RELEASE_FILE_PATH"`
	Owner           string `envconfig:"OWNER"`
	RepoFullName    string `envconfig:"REPO_FULL_NAME"`
	Repo            string `envconfig:"REPO"`
	BaseBranch      string `envconfig:"BASE_BRANCH"`
	NewTag          string `envconfig:"NEW_TAG"`
	UserName        string `envconfig:"USER_NAME"`
	UserEmail       string `envconfig:"USER_EMAIL"`
}

func main() {
	ctx := context.Background()
	var e env
	if err := envconfig.Process("INPUT", &e); err != nil {
		log.Fatal(err.Error())
	}

	user := gitConfig{userName: e.UserName, userEmail: e.UserEmail}
	gitCli := newGitClient(ctx, e.GithubToken, e.RepoFullName, user)
	oldTag, newNode, yamlPath, parseFile := readReleaseFile(gitCli.file, e.ReleaseFilePath)
	newNode, newTag := generateTag(newNode, oldTag, e.NewTag)
	log.Printf("tag: %v", newTag)
	branch := gitCli.Checkout(newTag)
	writeFile(yamlPath, gitCli.file, parseFile, newNode, e.ReleaseFilePath)
	gitCli.Commit(e.ReleaseFilePath, newTag)
	gitCli.Push(ctx, e.Owner)
	ghCli := newGHClient(e.GithubToken)
	ghCli.newPullRequest(ctx, newTag, e.BaseBranch, e.Repo, e.Owner, branch)
}
