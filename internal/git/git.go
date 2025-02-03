// TODO: review
package git

import (
	"fmt"
	"os"

	"github.com/CervinoB/scannercli/cmd/state"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/spf13/viper"
)

type Git struct {
	gs       *state.GlobalState
	RepoPath string
	TagName  string
}

// CloneRepo clones a Git repository from the specified URL to the given path.
// It logs the cloning process and returns an error if the cloning fails.
//
// Parameters:
//   - url: The URL of the Git repository to clone.
//   - path: The local file system path where the repository should be cloned.
//
// Returns:
//   - error: An error object if the cloning process fails, otherwise nil.
func CloneRepo(gs *state.GlobalState, url string, path string) error {
	// fmt.Println("Clonando repositório:", url)
	gs.Logger.Infof("Clonando repositório: %s", url)
	gs.Logger.Debug("Path do repositório:", path)
	_, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})
	return err
}

func CheckoutTag(gs *state.GlobalState, tagName string) error {
	repo, err := git.PlainOpen(viper.GetString("path"))
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
	gs.Logger.Info("git show-ref --head HEAD")
	ref, err := repo.Head()
	if err != nil {
		return err
	}
	fmt.Println(ref.Hash())
	return err
}

func Fetch(gs *state.GlobalState) error {
	repo, err := git.PlainOpen(viper.GetString("path"))
	if err != nil {
		gs.Logger.Error("Error fetching:", err)
		return err
	}

	return fetch(gs, repo)
}

func fetch(gs *state.GlobalState, repo *git.Repository) error {
	origin, err := repo.Remote("origin")
	if err != nil {
		return err
	}
	gs.Logger.Info("git fetch " + origin.Config().Name + " " + origin.Config().URLs[0])
	err = repo.Fetch(&git.FetchOptions{
		RemoteName: "origin",
		Progress:   os.Stdout,
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return err
	}

	// Check if there are updates
	if err == nil {
		gs.Logger.Info("Updates found, pulling changes")
		worktree, err := repo.Worktree()
		if err != nil {
			return err
		}
		err = worktree.Pull(&git.PullOptions{
			RemoteName: "origin",
			Progress:   os.Stdout,
		})
		if err != nil && err != git.NoErrAlreadyUpToDate {
			return err
		}
	} else {
		gs.Logger.Info("No updates found")
	}

	return nil
}
