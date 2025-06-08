package git

import (
	"os"
	"testing"
)

func TestCloneRepository(t *testing.T) {
	// This is a placeholder for the actual test implementation.
	// In a real test, you would set up a mock repository or use a test repository URL.
	repoURL := "https://github.com/twentyhq/twenty.git"
	targetDir := "repo/twenty"

	// Clean up targetDir before running the test
	if err := os.RemoveAll(targetDir); err != nil {
		t.Fatalf("Failed to remove targetDir before test: %v", err)
	}

	err := CloneRepository(repoURL, targetDir)
	if err != nil {
		t.Fatalf("TestCloneRepository failed: %v", err)
	}

	// Check if the .git directory exists in the targetDir
	gitDir := targetDir + "/.git"
	if _, err := os.Stat(gitDir); err != nil {
		if os.IsNotExist(err) {
			t.Errorf(".git directory does not exist in cloned repository")
		} else {
			t.Errorf("Error checking .git directory: %v", err)
		}
	}
}

func TestPullLatestChanges(t *testing.T) {
	// This is a placeholder for the actual test implementation.
	// In a real test, you would set up a mock repository or use a test repository URL.
	repoDir := "repo/twenty"

	// Clean up repoDir before running the test
	if err := os.RemoveAll(repoDir); err != nil {
		t.Fatalf("Failed to remove repoDir before test: %v", err)
	}

	err := CloneRepository("https://github.com/twentyhq/twenty.git", repoDir)
	if err != nil {
		t.Fatalf("TestPullLatestChanges failed: %v", err)
	}

	err = PullLatestChanges(repoDir)
	if err != nil {
		t.Fatalf("TestPullLatestChanges failed: %v", err)
	}
}
