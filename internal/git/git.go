package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// CloneRepository clones a Git repository from the given URL into the specified target directory.
func CloneRepository(repoURL, targetDir string) error {

	cmd := exec.Command("git", "clone", repoURL, targetDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}
	return nil
}

// ListBranches lists all branches in the specified Git repository directory.
func ListBranches(repoDir string) ([]string, error) {

	cmd := exec.Command("git", "-C", repoDir, "branch", "--list")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list branches: %w", err)
	}

	branches := strings.Fields(string(output))
	return branches, nil
}

// PullLatestChanges pulls the latest changes from the remote repository.
func PullLatestChanges(repoDir string) error {

	cmd := exec.Command("git", "-C", repoDir, "pull")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to pull latest changes: %w", err)
	}
	return nil
}

// CheckoutTag checks out a specific tag in the given repository directory.
func CheckoutTag(repoDir, tag string) error {

	cmd := exec.Command("git", "-C", repoDir, "checkout", tag)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to checkout tag %s: %w", tag, err)
	}
	return nil
}
