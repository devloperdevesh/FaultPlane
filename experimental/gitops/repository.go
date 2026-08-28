package gitops

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	ErrRepositoryPathRequired = errors.New("gitops repository path is required")
	ErrNotRepository          = errors.New("path is not a git repository")
	ErrGitUnavailable         = errors.New("git executable is unavailable")
)

type Repository struct {
	Path string
}

type Status struct {
	Path       string    `json:"path"`
	Branch     string    `json:"branch"`
	Commit     string    `json:"commit"`
	Remote     string    `json:"remote,omitempty"`
	Clean      bool      `json:"clean"`
	ObservedAt time.Time `json:"observed_at"`
}

func NewRepository(path string) (*Repository, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil, ErrRepositoryPathRequired
	}

	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("stat repository path: %w", err)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("%w: path is not a directory", ErrNotRepository)
	}

	return &Repository{
		Path: path,
	}, nil
}

func (r *Repository) Inspect(ctx context.Context) (Status, error) {
	if r == nil || strings.TrimSpace(r.Path) == "" {
		return Status{}, ErrRepositoryPathRequired
	}

	if _, err := os.Stat(r.Path); err != nil {
		return Status{}, fmt.Errorf("stat repository: %w", err)
	}

	if _, err := r.git(ctx, "rev-parse", "--is-inside-work-tree"); err != nil {
		return Status{}, fmt.Errorf("%w: %s", ErrNotRepository, strings.TrimSpace(err.Error()))
	}

	branch, err := r.git(ctx, "branch", "--show-current")
	if err != nil {
		return Status{}, fmt.Errorf("read git branch: %w", err)
	}

	commit, err := r.git(ctx, "rev-parse", "HEAD")
	if err != nil {
		return Status{}, fmt.Errorf("read git commit: %w", err)
	}

	remote, err := r.git(ctx, "remote", "get-url", "origin")
	if err != nil {
		remote = ""
	}

	porcelain, err := r.git(ctx, "status", "--porcelain")
	if err != nil {
		return Status{}, fmt.Errorf("read git status: %w", err)
	}

	return Status{
		Path:       r.Path,
		Branch:     strings.TrimSpace(branch),
		Commit:     strings.TrimSpace(commit),
		Remote:     strings.TrimSpace(remote),
		Clean:      strings.TrimSpace(porcelain) == "",
		ObservedAt: time.Now().UTC(),
	}, nil
}

func (r *Repository) git(ctx context.Context, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = r.Path

	output, err := cmd.CombinedOutput()
	if err != nil {
		if errors.Is(ctx.Err(), context.Canceled) {
			return "", context.Canceled
		}

		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return "", context.DeadlineExceeded
		}

		if _, lookupErr := exec.LookPath("git"); lookupErr != nil {
			return "", ErrGitUnavailable
		}

		message := strings.TrimSpace(string(output))
		if message == "" {
			message = err.Error()
		}

		return "", errors.New(message)
	}

	return string(output), nil
}
