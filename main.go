package main

import (
	"context"
	"log"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

type env struct {
	GithubToken     string `split_words:"true"`
	ReleaseFilePath string `split_words:"true"`
	RepoFullName    string `split_words:"true"`
	BaseBranch      string `split_words:"true"`
	UserName        string `split_words:"true"`
	UserEmail       string `split_words:"true"`
	NextSemverLevel string `split_words:"true"`
}

const (
	repoFullNameDelimiter = "/"
)

func main() {
	ctx := context.Background()
	var e env
	if err := envconfig.Process("INPUT", &e); err != nil {
		log.Fatal(err.Error())
	}

	separatedRepoFullName := strings.Split(e.RepoFullName, repoFullNameDelimiter)
	ownerName, repositoryName := separatedRepoFullName[0], separatedRepoFullName[1]

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
	gitCli.Push(ctx, ownerName)
	ghCli := newGHClient(e.GithubToken)
	ghCli.newPullRequest(ctx, newTag, e.BaseBranch, repositoryName, ownerName, branch)
}
