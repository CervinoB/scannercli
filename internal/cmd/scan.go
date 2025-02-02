/*
Copyright © 2025 Joao Cervino jcervinobarbosa@gmail.com
*/
package cmd

import (
	"context"
	"errors"

	"github.com/CervinoB/scannercli/cmd/state"
	"github.com/CervinoB/scannercli/internal/ui/pb"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type scanCmd struct {
	gs *state.GlobalState
}

func (c *scanCmd) run(cmd *cobra.Command, args []string) error {
	pbs := createProgressBars(10, 10, 5)
	var l logrus.FieldLogger = c.gs.Logger
	defer func(){
		if err == nil {
			l.Debug("Everything has finished, exiting scannercli normally!")
		}		else {
			l.WithError(err).Debug("Everything has finished, exiting scannercli with an error!")
		}
	}()

	printBanner(c.gs)

	globalCtx ,globalCanel := context.WithCancel(c.gs.Ctx)
	defer globalCanel()

	emitEvent := func(evt *event.Event) func() {
		waitDone := c.gs.Events.Emit(evt)
		return func() {
			waitCtx, waitCancel := context.WithTimeout(globalCtx, waitEventDoneTimeout)
			defer waitCancel()
			if werr := waitDone(waitCtx); werr != nil {
				logger.WithError(werr).Warn()
			}
		}
	}
	
		defer func() {
		waitExitDone := emitEvent(&event.Event{
			Type: event.Exit,
			Data: &event.ExitData{Error: err},
		})
		waitExitDone()
		c.gs.Events.UnsubscribeAll()
	}()

	if err = c.setupTracerProvider(globalCtx, test); err != nil {
		return err
	}
	

	initBar := pb.New(pb.WithConstLeft("Init"))
	l.Info("Starting repository scan...")

	initBar.Modify(pb.WithConstProgress(0, "Starting outputs"))
	l.Infof("Starting scan for: %s", args)

	// if docker {
	// 	runDockerizedScan(target)
	// } else {
	// 	runLocalScan(target)
	// }

	// 2. Executar scanners (SonarQube, ESLint)

	// 3. Gerar relatório

	l.Info("Scan completed successfully.")
	return nil
}

func runDockerizedScan(target string) {
	logrus.Info("Running dockerized scanners")
	// TODO: Implementar lógica de Docker
}

func runLocalScan(target string) {
	logrus.Info("Running local scanners")
	// TODO: Chamar módulos de análise
}

func createProgressBars(num, padding, colIdx int) []*pb.ProgressBar {
	pbs := make([]*pb.ProgressBar, num)
	for i := 0; i < num; i++ {
		left := fmt.Sprintf("left %d", i)
		rightCol1 := fmt.Sprintf("right %d", i)
		progress := 0.0
		status := pb.Running
		if i == colIdx {
			pad := strings.Repeat("+", padding)
			left += pad
			rightCol1 += pad
			progress = 1.0
			status = pb.Done
		}
		pbs[i] = pb.New(
			pb.WithLeft(func() string { return left }),
			pb.WithStatus(status),
			pb.WithProgress(func() (float64, []string) {
				return progress, []string{rightCol1, "000"}
			}),
		)
	}
	r

func getCmdScan(gs *state.GlobalState) *cobra.Command {
	s := &scanCmd{gs: gs}
	scanCmd := &cobra.Command{
		Use:   "scan",
		Short: "Scan a repository for code smells",
		Long: `Scan a repository using SonarQube, ESLint, and other tools to detect code smells.

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Args: exactArgsWithMsg(1, "Repository URL to scan"),
		RunE: s.run,
	}

	scanCmd.Flags().StringP("repo", "t", "", "Repository URL to scan")
	scanCmd.Flags().StringP("commit", "c", "HEAD", "Commit hash or tag to analyze")
	scanCmd.Flags().StringP("scanner", "s", "", "Scanner to use (sonarqube, eslint, etc.)")
	scanCmd.Flags().StringP("config", "f", "", "Configuration file to use")

	// rootCmd.AddCommand(scanCmd)
	return scanCmd
}

// scanCmd represents the scan command
// var scanCmdOld = &cobra.Command{
// 	Use:   "scan [target]",
// 	Short: "Scan a repository for code smells",
// 	Long: `Scan a repository using SonarQube, ESLint, and other tools to detect code smells.

// Cobra is a CLI library for Go that empowers applications.
// This application is a tool to generate the needed files
// to quickly create a Cobra application.`,
// 	Args: cobra.ExactArgs(1),
// 	Run: func(cmd *cobra.Command, args []string) {
// 		fmt.Println("scan called")
// 		logrus.Info("Starting repository scan...")

// 		// TODO: Implementar lógica de scan
// 		// 1. Carregar configurações
// 		target := args[0]
// 		logrus.Infof("Starting scan for: %s", target)
// 		// if docker {
// 		// 	runDockerizedScan(target)
// 		// } else {
// 		// 	runLocalScan(target)
// 		// }

// 		// 2. Executar scanners (SonarQube, ESLint)

// 		// 3. Gerar relatório

// 		logrus.Info("Scan completed successfully.")
// 	},
// }
