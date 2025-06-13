package api

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/CervinoB/scannercli/internal/logging"
)

func ExecSonarScanner(projectKey, token, sonarHost, sourcePath string, debug bool) error {
	cmd := exec.Command(
		"sonar-scanner",
		"-Dsonar.projectKey="+projectKey,
		"-Dsonar.sources=.",
		"-Dsonar.host.url="+sonarHost,
		"-Dsonar.token="+token,
	)

	cmd.Dir = filepath.Clean(sourcePath)

	if debug {
		logging.Logger.Debugf("Running sonar-scanner in %s with command: %s", sourcePath, cmd.String())

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		logging.Logger.Debugf("Stdout: %v", cmd.Stdout)
		logging.Logger.Debugf("Stderr: %v", cmd.Stderr)
	} else {
		var outBuf, errBuf bytes.Buffer
		cmd.Stdout = &outBuf
		cmd.Stderr = &errBuf
	}

	logging.Logger.Printf("Running sonar-scanner in %s", sourcePath)

	done := make(chan struct{})
	go func() {
		chars := []rune{'|', '/', '-', '\\'}
		i := 0
		for {
			select {
			case <-done:
				fmt.Print("\r") // Clear loader
				return
			default:
				fmt.Printf("\rScanning... %c", chars[i%len(chars)])
				time.Sleep(200 * time.Millisecond)
				i++
			}
		}
	}()

	err := cmd.Run()
	close(done)
	fmt.Print("\r")

	return err
}
