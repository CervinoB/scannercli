package sonar

import (
	"github.com/CervinoB/scannercli/cmd/state"
)

type Scanner struct {
	gs *state.GlobalState
	// TODO: adicionar outros parâmetros

}

func (s *Scanner) New(scanner string, args []string) error {
	l := s.gs.Logger

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
