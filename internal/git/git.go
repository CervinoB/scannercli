package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// CloneRepository clones a Git repository from the given URL into the specified target directory.
func CloneRepository(repoURL, targetDir string) error {
	// Check if the target directory exists
	if _, err := os.Stat(targetDir); err == nil {
		// Directory exists, ensure we are on main branch before pulling
		checkoutCmd := exec.Command("git", "-C", targetDir, "checkout", "main")
		checkoutCmd.Stdout = os.Stdout
		checkoutCmd.Stderr = os.Stderr
		if err := checkoutCmd.Run(); err != nil {
			return fmt.Errorf("failed to checkout main branch: %w", err)
		}
		// Pull latest changes
		if err := PullLatestChanges(targetDir); err != nil {
			return fmt.Errorf("repository exists but failed to pull latest changes: %w", err)
		}
		return nil
	}

	// Directory does not exist, clone the repository
	cmd := exec.Command("git", "clone", repoURL, targetDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}
	return nil
}

// ListTags lists all tags in the specified Git repository directory.
func ListTags(repoDir string) ([]string, error) {

	cmd := exec.Command("git", "-C", repoDir, "tag", "--list")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}

	tags := strings.Fields(string(output))
	return tags, nil
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
