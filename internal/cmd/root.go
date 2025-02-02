/*
Copyright © 2025 Joao Cervino jcervinobarbosa@gmail.com
*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"

	"github.com/CervinoB/sonarcli/cmd/state"
	"github.com/CervinoB/sonarcli/lib/consts"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.k6.io/k6/errext"
	"go.k6.io/k6/errext/exitcodes"
)

type rootCommand struct {
	debug bool
	cmd   *cobra.Command
	State *state.GlobalState
}

func newRootCommand(gs *state.GlobalState) *rootCommand {
	c := &rootCommand{
		State: gs,
	}

	// the base command when called without any subcommands.
	rootCmd := &cobra.Command{
		Use:   "sonarcli",
		Short: "a next-generation load generator",
		Long: "\n" + `SonarCLI facilita a análise de código estático utilizando o SonarScanner e 
	a coleta de métricas de code smell ao longo do tempo` + "\n" + consts.Banner(),
		// SilenceUsage:      true,
		// SilenceErrors:     true,
		PersistentPreRunE: c.persistentPreRunE,
		Version:           versionString(),
	}

	rootCmd.SetVersionTemplate(
		`{{with .Name}}{{printf "%s " .}}{{end}}{{printf "v%s\n" .Version}}`,
	)
	rootCmd.PersistentFlags().BoolVar(&c.debug, "debug", false, "enable debug mode")
	rootCmd.PersistentFlags().StringVar(&c.State.CfgFile, "config", "", "config file (default is $HOME/.sonarcli.yaml)")
	rootCmd.PersistentFlags().BoolVar(&c.State.Docker, "docker", false, "run scanners in docker containers")

	subCommands := []func(*state.GlobalState) *cobra.Command{
		getCmdArchive, getCmdCloud, getCmdNewScript, getCmdInspect,
		getCmdLogin, getCmdPause, getCmdResume, getCmdScale, getCmdRun,
		getCmdStats, getCmdStatus, getCmdVersion,
	}

	for _, sc := range subCommands {
		rootCmd.AddCommand(sc(gs))
	}

	c.cmd = rootCmd
	return c
}

func (c *rootCommand) persistentPreRunE(_ *cobra.Command, _ []string) error {
	c.initLogger()
	c.State.Logger.Debugf("k6 version: v%s", fullVersion())
	return nil
}

func (c *rootCommand) execute() {
	ctx, cancel := context.WithCancel(c.State.Ctx)
	c.State.Ctx = ctx

	exitCode := -1
	defer func() {
		cancel()
		c.State.OSExit(exitCode)
	}()

	defer func() {
		if r := recover(); r != nil {
			exitCode = int(exitcodes.GoPanic)
			err := fmt.Errorf("unexpected sonarcli panic: %s\n%s", r, debug.Stack())
			c.State.Logger.Error(err)
		}
	}()

	err := c.cmd.Execute()
	if err == nil {
		exitCode = 0
		return
	}

	var ecerr errext.HasExitCode
	if errors.As(err, &ecerr) {
		exitCode = int(ecerr.ExitCode())
	}

	errText, fields := errext.Format(err)
	c.State.Logger.WithFields(fields).Error(errText)
}

func Execute() {
	gs := state.NewState(context.Background())
	newRootCommand(gs).execute()
}

func (c *rootCommand) initLogger() {
	if c.debug {
		c.State.Logger.SetLevel(logrus.DebugLevel)
		c.State.Logger.Debug("Debug mode enabled")
	}
}
