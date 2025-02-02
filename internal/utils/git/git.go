// TODO: review
package git

import (
	"fmt"
	"os"
	. "sonarcli/internal/utils/common"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// CloneRepo clones a Git repository from the specified URL to the given path.
// It logs the cloning process and returns an error if the cloning fails.
//
// Parameters:
//   - url: The URL of the Git repository to clone.
//   - path: The local file system path where the repository should be cloned.
//
// Returns:
//   - error: An error object if the cloning process fails, otherwise nil.
func CloneRepo(url, path string) error {
	// fmt.Println("Clonando repositório:", url)
	Info("Clonando repositório: %s", url)
	_, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})
	return err
}

func CheckoutTag(repoPath, tagName string) error {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	fmt.Println("Alternando para tag:", tagName)
	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName("refs/tags/" + tagName),
	})
	Info("git show-ref --head HEAD")
	ref, err := repo.Head()
	CheckIfError(err)
	fmt.Println(ref.Hash())
	return err
}
