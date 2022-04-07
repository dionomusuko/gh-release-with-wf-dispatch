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
	userName   string
	token      string
}

func newGitClient(ctx context.Context, token, userName, repository string) *gitClient {
	fs := memfs.New()

	// https://<GITHUB_ACCESS_TOKEN>@github.com/<OWNEr>/<REPO>.git
	repo, err := git.CloneContext(ctx, memory.NewStorage(), fs, &git.CloneOptions{
		URL: fmt.Sprintf("https://%s@github.com/%s/%s.git", token, userName, repository),
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
		userName:   userName,
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
		log.Fatalf("falied to chckout repository")
	}
	return branch.String()
}

func (g *gitClient) Commit(ctx context.Context, filePath, newTag string) {
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

func (g *gitClient) Push(ctx context.Context) {
	if err := g.repository.PushContext(ctx, &git.PushOptions{
		Progress: os.Stdout,
		Auth: &http.BasicAuth{
			Username: g.userName,
			Password: g.token,
		},
		RemoteName: "origin",
	}); err != nil {
		log.Fatalf("failed to push: %v", err)
	}

}
