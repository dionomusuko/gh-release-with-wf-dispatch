package main

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type ghClient struct {
	client *github.Client
}

func newGHClient(ctx context.Context, token string) *ghClient {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return &ghClient{
		client: client,
	}
}

func (g *ghClient) createPullRequest(
	ctx context.Context, nextTag, baseBranch, repo, owner, branchName string,
) (string, error) {
	newPR := &github.NewPullRequest{
		Title: github.String("chore: release " + nextTag),
		Head:  github.String(owner + ":" + branchName),
		Base:  github.String(baseBranch),
		Body:  github.String("Release " + nextTag),
	}
	pr, _, err := g.client.PullRequests.Create(ctx, owner, repo, newPR)
	if err != nil {
		fmt.Println("failed to create pull request")
		return "", err
	}
	return pr.GetHTMLURL(), nil
}
