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
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
)

type gitClient struct {
	repository *git.Repository
	worktree   *git.Worktree
	file       billy.Filesystem
	token      string
}

func newGitClient(ctx context.Context, token, repository string) *gitClient {
	fs := memfs.New()

	// https://<GITHUB_TOKEN>@github.com/<REPO>.git
	repo, err := git.CloneContext(ctx, memory.NewStorage(), fs, &git.CloneOptions{
		URL: fmt.Sprintf("https://%s@github.com/%s.git", token, repository),
	})
	if err != nil {
		log.Fatalf("falied to clone repository")
	}
	w, err := repo.Worktree()
	if err != nil {
		log.Fatalf("falied to get worktree")
	}
	return &gitClient{
		repository: repo,
		worktree:   w,
		file:       fs,
		token:      token,
	}
}

// Checkout new branch
func (g *gitClient) Checkout(newTag string) string {
	branch := plumbing.ReferenceName("refs/heads/release-" + newTag)
	if err := g.worktree.Checkout(&git.CheckoutOptions{
		Create: true,
		Branch: branch,
	}); err != nil {
		log.Fatalf("falied to chckout repository: %v", err)
	}
	return branch.String()
}

func (g *gitClient) Commit(filePath, newTag string) {
	// Create Commit
	if _, err := g.worktree.Add(filePath); err != nil {
		log.Fatalf("falied to add")
	}
	_, err := g.worktree.Commit("chore: release-"+newTag, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "github-actions",
			Email: "github-actions@github.com",
			When:  time.Now(),
		}})
	if err != nil {
		log.Fatalf("falied to commit")
	}
}

func (g *gitClient) Push(ctx context.Context, owner string) {
	if err := g.repository.PushContext(ctx, &git.PushOptions{
		Progress: os.Stdout,
		Auth: &http.BasicAuth{
			Username: owner,
			Password: g.token,
		},
		RemoteName: "origin",
	}); err != nil {
		log.Fatalf("failed to push: %v", err)
	}

}
