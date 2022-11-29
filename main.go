package main

import (
	"context"
	"fmt"
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
	Assignees       string `split_words:"true"`
}

const (
	repoFullNameDelimiter = "/"
	inputPrefix           = "INPUT"
)

func main() {
	ctx := context.Background()
	var e env
	if err := envconfig.Process(inputPrefix, &e); err != nil {
		fmt.Printf("failed to load inputs: %s\n", err.Error())
		panic(err.Error())
	}

	ownerName, repositoryName := splitRepoFullName(e.RepoFullName)
	gitConf := gitConfig{userName: e.UserName, userEmail: e.UserEmail}
	gitCli, err := newGitClient(ctx, e.GithubToken, e.RepoFullName, gitConf)
	if err != nil {
		fmt.Printf("failed to create git client: %s\n", err.Error())
		panic(err.Error())
	}
	currentTag, newNode, yamlPath, parseFile, err := readReleaseFile(gitCli.file, e.ReleaseFilePath)
	if err != nil {
		fmt.Printf("failed to read release file: %s\n", err.Error())
		panic(err.Error())
	}
	nextTag, err := newSemver(currentTag, e.NextSemverLevel)
	if err != nil {
		fmt.Printf("failed to parse semver for %s: %s", currentTag, err.Error())
		panic(err.Error())
	}
	newNode.Value = nextTag
	branchName, err := gitCli.Checkout("refs/heads/release-" + nextTag)
	if err != nil {
		fmt.Printf("failed to checkout %s\n", branchName)
		panic(err.Error())
	}
	if err := writeFile(yamlPath, gitCli.file, parseFile, newNode, e.ReleaseFilePath); err != nil {
		fmt.Println("failed to write file")
		panic(err.Error())
	}

	// Git operation
	if err := gitCli.Add(e.ReleaseFilePath); err != nil {
		panic(err.Error())
	}
	if err := gitCli.Commit("chore: release " + nextTag); err != nil {
		panic(err.Error())
	}
	if err := gitCli.Push(ctx); err != nil {
		panic(err.Error())
	}

	// GitHub operation
	ghCli := newGHClient(ctx, e.GithubToken, ownerName, repositoryName)
	pr, err := ghCli.createPullRequest(ctx, nextTag, e.BaseBranch, branchName)
	if err != nil {
		fmt.Println("failed to create pull request")
		panic(err.Error())
	}
	fmt.Println(pr.GetHTMLURL())
	// TODO: convert e.Assignees to []string
	if err := ghCli.addAssignees(ctx, *pr.Number, []string{e.Assignees}); err != nil {
		fmt.Println("failed to add assignees")
		panic(err.Error())
	}
}

func splitRepoFullName(fullName string) (string, string) {
	separatedRepoFullName := strings.Split(fullName, repoFullNameDelimiter)
	return separatedRepoFullName[0], separatedRepoFullName[1]
}
