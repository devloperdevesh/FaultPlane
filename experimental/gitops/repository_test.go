package gitops

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestInspectRealRepository(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	root := t.TempDir()

	runGit(t, root, "init")
	runGit(t, root, "config", "user.email", "test@example.com")
	runGit(t, root, "config", "user.name", "FaultPlane Test")

	file := filepath.Join(root, "state.txt")
	if err := os.WriteFile(file, []byte("real gitops state\n"), 0o644); err != nil {
		t.Fatalf("write test file: %v", err)
	}

	runGit(t, root, "add", "state.txt")
	runGit(t, root, "commit", "-m", "initial state")

	repo, err := NewRepository(root)
	if err != nil {
		t.Fatalf("NewRepository() error = %v", err)
	}

	status, err := repo.Inspect(ctx)
	if err != nil {
		t.Fatalf("Inspect() error = %v", err)
	}

	if status.Path != root {
		t.Fatalf("Path = %q, want %q", status.Path, root)
	}

	if status.Commit == "" {
		t.Fatal("Commit is empty")
	}

	if !status.Clean {
		t.Fatal("expected clean repository")
	}

	if status.ObservedAt.IsZero() {
		t.Fatal("ObservedAt is zero")
	}
}

func TestInspectDetectsDirtyRepository(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	root := t.TempDir()

	runGit(t, root, "init")
	runGit(t, root, "config", "user.email", "test@example.com")
	runGit(t, root, "config", "user.name", "FaultPlane Test")

	file := filepath.Join(root, "state.txt")
	if err := os.WriteFile(file, []byte("initial\n"), 0o644); err != nil {
		t.Fatalf("write initial file: %v", err)
	}

	runGit(t, root, "add", "state.txt")
	runGit(t, root, "commit", "-m", "initial state")

	if err := os.WriteFile(file, []byte("changed\n"), 0o644); err != nil {
		t.Fatalf("modify file: %v", err)
	}

	repo, err := NewRepository(root)
	if err != nil {
		t.Fatalf("NewRepository() error = %v", err)
	}

	status, err := repo.Inspect(ctx)
	if err != nil {
		t.Fatalf("Inspect() error = %v", err)
	}

	if status.Clean {
		t.Fatal("expected dirty repository")
	}
}

func TestNewRepositoryRequiresPath(t *testing.T) {
	if _, err := NewRepository(" "); err != ErrRepositoryPathRequired {
		t.Fatalf("error = %v, want %v", err, ErrRepositoryPathRequired)
	}
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()

	cmd := exec.Command("git", args...)
	cmd.Dir = dir

	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, output)
	}
}
