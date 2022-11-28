package main

import (
	"context"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type ghClient struct {
	client *github.Client
	owner  string
	repo   string
}

func newGHClient(ctx context.Context, token string, owner, repo string) *ghClient {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return &ghClient{
		client: client,
		owner:  owner,
		repo:   repo,
	}
}

func (g *ghClient) createPullRequest(
	ctx context.Context, nextTag, baseBranch, branchName string,
) (*github.PullRequest, error) {
	newPR := &github.NewPullRequest{
		Title: github.String("chore: release " + nextTag),
		Head:  github.String(g.owner + ":" + branchName),
		Base:  github.String(baseBranch),
		Body:  github.String("Release " + nextTag),
	}
	pr, _, err := g.client.PullRequests.Create(ctx, g.owner, g.repo, newPR)
	if err != nil {
		return nil, err
	}
	return pr, nil
}

func (g *ghClient) addAssignees(ctx context.Context, number int, assignees []string) error {
	_, _, err := g.client.Issues.AddAssignees(ctx, g.owner, g.repo, number, assignees)
	if err != nil {
		return err
	}
	return nil
}
