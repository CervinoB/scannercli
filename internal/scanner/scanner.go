package scanner

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/CervinoB/scannercli/cmd/state"
)

type Scanner struct {
	// TODO: adicionar outros parâmetros

}

const (
	baseURL     = "http://localhost:9000"
	loginPath   = "/api/authentication/login"
	createPath  = "/api/projects/create"
	username    = "admin"
	password    = "zy3fnVnvKLw4dca!"
	projectKey  = "test222"
	projectName = "test222"
	mainBranch  = "main"
)

func New(gs *state.GlobalState, scanner string, project string) error {
	var l = gs.Logger

	l.Infof("Running scanner: %s", scanner)

	// scannerPath := filepath.Join(s.gs.CfgFile.SonarScanner.Path, scanner)
	// cmd := exec.Command(scannerPath, args...)
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	// err := cmd.Run()
	// if err != nil {
	// 	return fmt.Errorf("error running scanner: %s", err)
	// }
	return nil
}

func ExecScanner(gs *state.GlobalState, scanner string, path string) error {
	l := gs.Logger
	l.Infof("Executing scanner: %s", scanner)

	// run scanner command
	// For example, if using SonarQube:
	cmd := exec.Command("sonar-scanner", "-Dsonar.projectKey=test", "-Dsonar.sources=.", "-Dsonar.host.url=http://127.0.0.1:9000", "-Dsonar.login=sqp_00a92e640d2cdc01c9334a466b0b850efcd3cf1b")
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error executing scanner: %w", err)
	}
	l.Debugf("Scanner executed successfully: %s", scanner)

	return nil
}
