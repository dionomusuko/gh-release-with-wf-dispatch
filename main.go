package main

import (
	"context"
	"log"

	"github.com/kelseyhightower/envconfig"
)

// TODO 構成
// gh actions work flow dispatchで リリースバージョンが取れるのでハンドリングする
// リリースバージョンあり-> その値に置き換え
// なし -> パッチリリースするよう値をインクリメント
// RELEASE fileは yamlとして考える
type env struct {
	GithubToken     string `envconfig:"GITHUB_TOKEN"`
	ReleaseFilePath string `envconfig:"RELEASE_FILE_PATH"`
	NewTag          string `envconfig:"NEW_TAG"`
	Owner           string `envconfig:"OWNER"`
	Repo            string `envconfig:"REPO"`
	BaseBranch      string `envconfig:"BASE_BRANCH"`
}

func main() {
	ctx := context.Background()
	var e env
	err := envconfig.Process("INPUT", &e)
	if err != nil {
		log.Fatal(err.Error())
	}

	gitCli := newGitClient(ctx, e.GithubToken, e.Owner, e.Repo)
	oldTag, newNode, yamlPath, parseFile := readReleaseFile(gitCli.file, e.ReleaseFilePath)
	newNode, newTag := generateTag(newNode, oldTag, e.NewTag)
	log.Printf("tag: %v", newTag)
	writeFile(yamlPath, gitCli.file, parseFile, newNode, e.ReleaseFilePath)
	branch := gitCli.Checkout(newTag)
	gitCli.Commit(ctx, e.ReleaseFilePath, newTag)
	gitCli.Push(ctx)
	ghCli := newGHClient(e.GithubToken)
	ghCli.newPullRequest(ctx, newTag, e.BaseBranch, e.Repo, e.Owner, branch)
}
