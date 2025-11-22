package helpers

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"log"
	"time"
)

type MyGit struct {
	repo *git.Repository
	auth *http.BasicAuth
}

// NewMyGit initializes auth from environment and return a ready MyGit object
func NewMyGit(token string) (*MyGit, error) {

	auth := &http.BasicAuth{
		Username: "oauth2",
		Password: token,
	}

	return &MyGit{auth: auth}, nil
}

func (g *MyGit) CloneRepository(clonePath, repoURL string) error {
	repo, err := git.PlainClone(clonePath, false, &git.CloneOptions{
		URL:               repoURL,
		Auth:              g.auth,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})
	if err != nil {
		if errors.Is(err, git.ErrRepositoryAlreadyExists) {
			log.Println("Repo already exists, opening instead...")
			repo, err = git.PlainOpen(clonePath)
			if err != nil {
				return fmt.Errorf("failed to open existing repo: %w", err)
			}
		} else {
			return fmt.Errorf("failed to clone repo: %w", err)
		}
	}

	g.repo = repo
	return nil
}

func (g *MyGit) Reference(branchRefName string) (*plumbing.Reference, error) {
	return g.repo.Reference(plumbing.NewBranchReferenceName(branchRefName), true)
}

func (g *MyGit) CheckoutBranch(branchName string, create bool) error {
	if g.repo == nil {
		return errors.New("repo is nil")
	}
	w, err := g.repo.Worktree()
	if err != nil {
		return err
	}

	branch := plumbing.NewBranchReferenceName(branchName)

	return w.Checkout(&git.CheckoutOptions{
		Branch: branch,
		Create: create,
	})
}

func (g *MyGit) AddAndCommit() (*plumbing.Hash, error) {
	if g.repo == nil {
		return nil, errors.New("repo is nil")
	}

	w, err := g.repo.Worktree()
	if err != nil {
		return nil, err
	}

	_, err = w.Add(".")
	if err != nil {
		return nil, err
	}

	// Check if the working tree is clean
	status, err := w.Status()
	if err != nil {
		return nil, err
	}

	if status.IsClean() {
		// Nothing to commit → silently ignoring
		return nil, nil
	}

	commitMessage := fmt.Sprintf("update posts %v", time.Now())
	commitHash, err := w.Commit(commitMessage, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Your Name",
			Email: "your@email",
			When:  time.Now(),
		},
	})
	if err != nil {
		return nil, err
	}

	return &commitHash, nil
}

func (g *MyGit) Push(branchName string) error {
	if g.repo == nil {
		return errors.New("repo is nil")
	}

	return g.repo.Push(&git.PushOptions{
		Auth: g.auth,
		RefSpecs: []config.RefSpec{
			config.RefSpec("refs/heads/" + branchName + ":refs/heads/" + branchName),
		},
	})
}
