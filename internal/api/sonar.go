package api

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/CervinoB/scannercli/internal/logging"
)

func ExecSonarScanner(projectKey, token, sonarHost, sourcePath string) error {
	cmd := exec.Command(
		"sonar-scanner",
		"-Dsonar.projectKey="+projectKey,
		"-Dsonar.sources=.",
		"-Dsonar.host.url="+sonarHost,
		"-Dsonar.token="+token,
	)

	cmd.Dir = filepath.Clean(sourcePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	logging.Logger.Printf("Running sonar-scanner in %s", sourcePath)

	return cmd.Run()
}
