/*
Copyright © 2025 Joao Cervino jcervinobarbosa@gmail.com
*/
package cmd

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/CervinoB/scannercli/cmd/state"
	"github.com/CervinoB/scannercli/lib/consts"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type rootCommand struct {
	debug       bool
	cmd         *cobra.Command
	globalState *state.GlobalState
}

func newRootCommand(gs *state.GlobalState) *rootCommand {
	c := &rootCommand{
		globalState: gs,
	}

	// the base command when called without any subcommands.
	rootCmd := &cobra.Command{
		Use:   "scannercli",
		Short: "Code smell analysis tool for NestJS projects",
		Long: "\n" + `scannercli facilita a análise de código estático utilizando o SonarScanner e 
	a coleta de métricas de code smell ao longo do tempo` + "\n" + consts.Banner(),
		PersistentPreRunE: c.persistentPreRunE,
		Version:           versionString(),
	}

	rootCmd.SetVersionTemplate(
		`{{with .Name}}{{printf "%s " .}}{{end}}{{printf "v%s\n" .Version}}`,
	)
	rootCmd.PersistentFlags().BoolVar(&c.debug, "debug", false, "enable debug mode")
	rootCmd.PersistentFlags().StringVar(&c.globalState.CfgFile, "config", "", "config file (default is $HOME/.scannercli.yaml)")
	rootCmd.PersistentFlags().BoolVar(&c.globalState.Docker, "docker", false, "run scanners in docker containers")

	subCommands := []func(*state.GlobalState) *cobra.Command{getCmdVersion, getCmdScan}

	for _, sc := range subCommands {
		rootCmd.AddCommand(sc(gs))
	}

	c.cmd = rootCmd
	return c
}

func (c *rootCommand) persistentPreRunE(_ *cobra.Command, _ []string) error {
	c.initLogger()
	c.globalState.Logger.Debugf("scannercli version: v%s", fullVersion())
	return nil
}

func (c *rootCommand) execute() {
	ctx, cancel := context.WithCancel(c.globalState.Ctx)
	c.globalState.Ctx = ctx

	exitCode := -1
	defer func() {
		cancel()
		c.globalState.OSExit(exitCode)
	}()

	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("unexpected scannercli panic: %s\n%s", r, debug.Stack())
			c.globalState.Logger.Error(err)
		}
	}()

	err := c.cmd.Execute()
	if err == nil {
		exitCode = 0
		return
	}

	CheckIfError(c.globalState, err)
}

func Execute() {
	gs := state.NewState(context.Background())
	newRootCommand(gs).execute()
}

func (c *rootCommand) initLogger() {
	if c.debug {
		c.globalState.Logger.SetLevel(logrus.DebugLevel)
		c.globalState.Logger.Debug("Debug mode enabled")
	}
}
