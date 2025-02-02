package docker

import (
	"os/exec"
)

func RunDockerCommand(args ...string) error {
	cmd := exec.Command("docker", args...)

	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
