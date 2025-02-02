/*
Copyright © 2025 Joao Cervino jcervinobarbosa@gmail.com
*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/CervinoB/scannercli/cmd/state"
	"github.com/CervinoB/scannercli/lib/consts"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type rootCommand struct {
	cmd *cobra.Command
	gs  *state.GlobalState
}

func newRootCommand(gs *state.GlobalState) *rootCommand {
	c := &rootCommand{
		gs: gs,
	}

	// the base command when called without any subcommands.
	rootCmd := &cobra.Command{
		Use:   "scannercli",
		Short: "Code smell analysis tool for NestJS projects",
		Long: "\n" + `scannercli facilita a análise de código estático utilizando o SonarScanner e 
	a coleta de métricas de code smell ao longo do tempo` + "\n" + consts.Banner(),
		PersistentPreRunE: c.persistentPreRunE,
		Version:           versionString(),
		Example:           "scannercli scan",
	}

	rootCmd.SetVersionTemplate(
		`{{with .Name}}{{printf "%s " .}}{{end}}{{printf "v%s\n" .Version}}`,
	)
	rootCmd.PersistentFlags().StringVar(&c.gs.CfgFile, "config", "", "config file (default is $HOME/.scannercli.yaml)")
	rootCmd.PersistentFlags().BoolVar(&c.gs.Docker, "docker", false, "run scanners in docker containers")

	subCommands := []func(*state.GlobalState) *cobra.Command{getCmdVersion, getCmdScan}

	for _, sc := range subCommands {
		rootCmd.AddCommand(sc(gs))
	}

	c.cmd = rootCmd
	return c
}

func (c *rootCommand) persistentPreRunE(_ *cobra.Command, _ []string) error {
	c.initLogger()

	c.gs.Logger.Debugf("scannercli version: v%s", fullVersion())
	return nil
}

func (c *rootCommand) execute() {
	ctx, cancel := context.WithCancel(c.gs.Ctx)
	c.gs.Ctx = ctx

	exitCode := -1
	defer func() {
		cancel()
		c.gs.OSExit(exitCode)
	}()

	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("unexpected scannercli panic: %s\n%s", r, debug.Stack())
			c.gs.Logger.Error(err)
		}
	}()

	err := c.cmd.Execute()
	if err == nil {
		exitCode = 0
		return
	}

	CheckIfError(c.gs, err)
}

func Execute() {
	gs := state.NewState(context.Background())
	// p := tea.NewProgram() //TODO: newRootCommand(gs) -> initialModel()
	newRootCommand(gs).execute()
}

func (c *rootCommand) initConfig() {
	if c.gs.CfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(c.gs.CfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".cobra")
	}

	viper.AutomaticEnv()

	err := viper.ReadInConfig()

	notFound := &viper.ConfigFileNotFoundError{}
	switch {
	case err != nil && !errors.As(err, notFound):
		cobra.CheckErr(err)
	case err != nil && errors.As(err, notFound):
		// The config file is optional, we shouldn't exit when the config is not found
		break
	default:
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func (c *rootCommand) initLogger() {
	c.gs.Logger.SetLevel(logrus.DebugLevel)
	c.gs.Logger.Debug("Debug mode enabled")
	c.gs.Logger.SetLevel(logrus.InfoLevel)
}
