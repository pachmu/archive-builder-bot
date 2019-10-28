package main

import (
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"time"
)

type GitClient interface {
	Pull(branch string) error
	CommitAndPush(commitMessage string) error
}

type gitClient struct {
	url      string
	dir      string
	user     string
	email    string
	password string
}

func GetNewGitClient(url string, dir string, user string, email string, password string) GitClient {
	return &gitClient{
		url:      url,
		dir:      dir,
		user:     user,
		email:    email,
		password: password,
	}
}

func (g *gitClient) Pull(branch string) error {
	r, err := git.PlainClone(g.dir, false, &git.CloneOptions{
		URL:               g.url,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		Auth: &http.BasicAuth{
			Username: g.user,
			Password: g.password,
		},
	})
	if err != nil {
		return fmt.Errorf("ropo clone failed, got %s", err.Error())
	}
	w, err := r.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get repo worktree, got %s", err.Error())
	}

	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branch),
		Create: true,
	})
	if err != nil {
		return fmt.Errorf("failed to checkout to branch %s, got %s", branch, err.Error())
	}
	err = w.Pull(&git.PullOptions{
		RemoteName:    "origin",
		ReferenceName: plumbing.NewBranchReferenceName(branch),
		Auth: &http.BasicAuth{
			Username: g.user,
			Password: g.password,
		},
	})
	if err != nil {
		if err.Error() == "already up-to-date" {
			return nil
		}
		return fmt.Errorf("failed to pull branch %s from origin, got %s", branch, err.Error())
	}
	return nil
}

func (g *gitClient) CommitAndPush(commitMessage string) error {
	r, err := git.PlainOpen(g.dir)
	if err != nil {
		return fmt.Errorf("ropo clone failed, got %s", err.Error())
	}
	w, err := r.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get repo worktree, got %s", err.Error())
	}
	commit, err := w.Commit(commitMessage, &git.CommitOptions{
		All: true,
		Author: &object.Signature{
			Name:  g.user,
			Email: g.email,
			When:  time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to get commit, got %s", err.Error())
	}
	_, err = r.CommitObject(commit)
	if err != nil {
		return fmt.Errorf("failed to commit, got %s", err.Error())
	}
	err = r.Push(&git.PushOptions{
		Auth: &http.BasicAuth{
			Username: g.user,
			Password: g.password,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to push, got %s", err.Error())
	}
	return nil
}
