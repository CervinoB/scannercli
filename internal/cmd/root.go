/*
Copyright © 2025 Joao Cervino jcervinobarbosa@gmail.com
*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime/debug"

	"github.com/CervinoB/scannercli/cmd/state"
	"github.com/CervinoB/scannercli/lib/consts"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type rootCommand struct {
	cmd   *cobra.Command
	gs    *state.GlobalState
	debug bool
}

func Execute() {
	gs := state.NewState(context.Background())
	// p := tea.NewProgram() //TODO: newRootCommand(gs) -> initialModel()
	newRootCommand(gs).execute()
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
	rootCmd.PersistentFlags().StringVar(&c.gs.CfgFile, "config", "", "config file (default is ./scannercli.yaml)")
	rootCmd.PersistentFlags().BoolVar(&c.gs.Docker, "docker", false, "run scanners in docker containers")
	rootCmd.PersistentFlags().BoolVarP(&c.debug, "debug", "d", false, "debug mode")

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
	c.initConfig(c.gs)

	c.ensureSonarContainerRunning()

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

func (c *rootCommand) initConfig(gs *state.GlobalState) {
	if c.gs.CfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(c.gs.CfgFile)
	} else {
		// Find home directory.
		// Search config in home directory with name ".cobra" (without extension).
		viper.SetConfigName("scannercli.yaml")
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")

		err := viper.ReadInConfig() // Find and read the config file
		if err != nil {             // Handle errors reading the config file
			gs.Logger.Error(err)
			panic(fmt.Errorf("fatal error config file: %w", err))
		}
		configs := viper.AllSettings()
		// Log all configurations
		logConfigs(gs, configs)

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

// Helper function to log configurations recursively
func logConfigs(gs *state.GlobalState, configs map[string]interface{}) {
	logger := gs.Logger
	for key, value := range configs {
		switch v := value.(type) {
		case map[string]interface{}:
			// If the value is a nested map, log the key and recursively log its contents
			logger.Debugf("Key: %s", key)
			logConfigs(gs, v)
		case []interface{}:
			// If the value is a slice, log each item in the slice
			logger.Debugf("Key: %s", key)
			for i, item := range v {
				logger.Debugf("  Item %d: %v", i, item)
			}
		default:
			// Log simple key-value pairs
			logger.Debugf("Key: %s, Value: %v", key, value)
		}
	}
}

func (c *rootCommand) initLogger() {
	if c.debug {
		c.gs.Logger.SetLevel(logrus.DebugLevel)
		c.gs.Logger.Debug("Debug mode enabled")
	} else {
		c.gs.Logger.SetLevel(logrus.InfoLevel)
	}
}

func (c *rootCommand) ensureSonarContainerRunning() string {
	cmd := exec.Command("docker", "ps", "-q", "--filter", "name=sonarqube-scanner")
	output, _ := cmd.Output()
	if len(output) > 0 {
		fmt.Println("SonarQube container already running.")
		return "<container_id>" // container is running
	}
	return ""
}
