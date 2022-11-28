package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
)

type gitClient struct {
	repository *git.Repository
	worktree   *git.Worktree
	file       billy.Filesystem
	token      string
	config     gitConfig
}

type gitConfig struct {
	userName  string
	userEmail string
}

func newGitClient(ctx context.Context, token, repository string, config gitConfig) (*gitClient, error) {
	fs := memfs.New()

	// https://<GITHUB_TOKEN>@github.com/<REPO>.git
	repo, err := git.CloneContext(ctx, memory.NewStorage(), fs, &git.CloneOptions{
		URL: fmt.Sprintf("https://%s@github.com/%s.git", token, repository),
	})
	if err != nil {
		fmt.Printf("failed to clone repository %s\n", repository)
		return nil, err
	}
	w, err := repo.Worktree()
	if err != nil {
		fmt.Println("failed to get worktree")
		return nil, err
	}
	return &gitClient{
		repository: repo,
		worktree:   w,
		file:       fs,
		token:      token,
		config:     config,
	}, nil
}

// Checkout new branch
func (g *gitClient) Checkout(nextTag string) string {
	branch := plumbing.ReferenceName("refs/heads/release-" + nextTag)
	if err := g.worktree.Checkout(&git.CheckoutOptions{
		Create: true,
		Branch: branch,
	}); err != nil {
		log.Fatalf("falied to chckout repository: %v", err)
	}
	return branch.String()
}

func (g *gitClient) Add(filePath string) error {
	if _, err := g.worktree.Add(filePath); err != nil {
		fmt.Printf("failed to git add file %s\n", filePath)
		return err
	}
	return nil
}

func (g *gitClient) Commit(msg string) error {
	commitHash, err := g.worktree.Commit(
		msg,
		&git.CommitOptions{
			Author: &object.Signature{
				Name:  g.config.userName,
				Email: g.config.userEmail,
				When:  time.Now(),
			},
		},
	)
	if err != nil {
		fmt.Println("failed to git commit")
		return err
	}
	fmt.Printf("commit %s has been created\n", commitHash.String())
	return nil
}

func (g *gitClient) Push(ctx context.Context) error {
	if err := g.repository.PushContext(
		ctx,
		&git.PushOptions{
			Progress:   os.Stdout,
			RemoteName: "origin",
		},
	); err != nil {
		fmt.Println("failed to git push")
		return err
	}
	return nil
}
