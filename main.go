package main

import (
	"context"
	"log"

	"github.com/kelseyhightower/envconfig"
)

type env struct {
	GithubToken     string `split_words:"true"`
	ReleaseFilePath string `split_words:"true"`
	Owner           string `split_words:"true"`
	RepoFullName    string `split_words:"true"`
	Repo            string `split_words:"true"`
	BaseBranch      string `split_words:"true"`
	UserName        string `split_words:"true"`
	UserEmail       string `split_words:"true"`
	NextSemverLevel string `split_words:"true"`
}

func main() {
	ctx := context.Background()
	var e env
	if err := envconfig.Process("INPUT", &e); err != nil {
		log.Fatal(err.Error())
	}

	user := gitConfig{userName: e.UserName, userEmail: e.UserEmail}
	gitCli := newGitClient(ctx, e.GithubToken, e.RepoFullName, user)
	currentTag, newNode, yamlPath, parseFile := readReleaseFile(gitCli.file, e.ReleaseFilePath)
	nextTag, err := newSemver(currentTag, e.NextSemverLevel)
	if err != nil {
		log.Printf("currentTag: %s\n", currentTag)
		log.Fatalf("failed to parse semver: %s", err.Error())
	}
	newNode, newTag := generateTag(newNode, currentTag, nextTag)
	branch := gitCli.Checkout(newTag)
	writeFile(yamlPath, gitCli.file, parseFile, newNode, e.ReleaseFilePath)
	gitCli.Commit(e.ReleaseFilePath, newTag)
	gitCli.Push(ctx, e.Owner)
	ghCli := newGHClient(e.GithubToken)
	ghCli.newPullRequest(ctx, newTag, e.BaseBranch, e.Repo, e.Owner, branch)
}
