package main

import (
	"context"
	"log"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

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

func (g *ghClient) newPullRequest(ctx context.Context, newTag, baseBranch, repo, owner, branchName string) {
	pr := &github.NewPullRequest{
		Title: github.String("chore: release " + newTag),
		Head:  github.String(owner + ":" + branchName),
		Base:  github.String(baseBranch),
		Body:  github.String("Release " + newTag),
	}
	if _, _, err := g.client.PullRequests.Create(ctx, owner, repo, pr); err != nil {
		log.Fatalf("faield to create pull request: %v", err)
	}
}
