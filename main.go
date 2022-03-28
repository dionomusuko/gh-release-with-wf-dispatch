package main

import (
	"context"
	"log"
	"time"

	"golang.org/x/oauth2"

	"github.com/google/go-github/v43/github"
	"github.com/kelseyhightower/envconfig"
)

// TODO 構成
// gh actions work flow dispatchで リリースバージョンが取れるのでハンドリングする
// リリースバージョンあり-> その値に置き換え
// なし -> パッチリリースするよう値をインクリメント
// RELEASE fileは yamlとして考える

type env struct {
	GithubToken        string `envconfig:"GITHUB_TOKEN"`
	GithubOrganization string `envconfig:"GITHUB_ORGANIZATION"`
	RepositoryName     string `envconfig:"REPOSITORY_NAME"`
	ReleaseFilePath    string `envconfig:"RELEASE_FILE_PATH"`
	OldTag             string `envconfig:"OLD_TAG"`
	NewTag             string `envconfig:"NEW_TAG"`
}

const (
	jobTimeout = 10 * 60 * time.Second
)

func main() {
	var e env
	err := envconfig.Process("INPUT", &e)
	if err != nil {
		log.Fatal(err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), jobTimeout)
	defer cancel()
	client := newGHClient(e.GithubToken)

	if e.NewTag == "" {
		client.client.Organizations.GetPackage(ctx, e.GithubOrganization, "maven", e.RepositoryName)
		increment("タグを入れる")
		return
	}
	replacement(e.NewTag)
}

type ghClient struct {
	client *github.Client
}

func newGHClient(token string) *ghClient {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return &ghClient{
		client: client,
	}
}
