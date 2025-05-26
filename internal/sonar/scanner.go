package scanner

import (
	"github.com/CervinoB/scannercli/cmd/state"
)

type Scanner struct {
	// TODO: adicionar outros parâmetros

}

func (s *Scanner) New(gs *state.GlobalState, scanner string, args []string) error {
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

func (s *Scanner) ExecScanner(gs *state.GlobalState, scanner string, args []string) error {
	l := gs.Logger
	l.Infof("Executing scanner: %s", scanner)

	return nil
}
