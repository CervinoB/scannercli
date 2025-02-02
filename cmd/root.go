/*
Copyright © 2025 Joao Cervino jcervinobarbosa@gmail.com
*/
package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type rootCommand struct {
	cfgFile string
	debug   bool
	docker  bool
	cmd     *cobra.Command
	Logger  *logrus.Logger
}

func newRootCommand(lg *logrus.Logger) *rootCommand {
	c := &rootCommand{
		Logger: lg,
	}

	// the base command when called without any subcommands.
	rootCmd := &cobra.Command{
		Use:   "sonarcli",
		Short: "a next-generation load generator",
		Long: `SonarCLI facilita a análise de código estático utilizando o SonarScanner e 
	a coleta de métricas de code smell ao longo do tempo.

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		// SilenceUsage:      true,
		// SilenceErrors:     true,
		PersistentPreRunE: c.persistentPreRunE,
		// Version:           versionString(),
	}

	rootCmd.SetVersionTemplate(
		`{{with .Name}}{{printf "%s " .}}{{end}}{{printf "v%s\n" .Version}}`,
	)
	rootCmd.PersistentFlags().StringVar(&c.cfgFile, "config", "", "config file (default is $HOME/.sonarcli.yaml)")
	rootCmd.PersistentFlags().BoolVar(&c.debug, "debug", false, "enable debug mode")
	rootCmd.PersistentFlags().BoolVar(&c.docker, "docker", false, "run scanners in docker containers")

	c.cmd = rootCmd
	return c
}
func (c *rootCommand) persistentPreRunE(_ *cobra.Command, _ []string) error {
	err := c.initLogger(c.Logger)
	if err != nil {
		return err
	}
	c.globalState.Logger.Debugf("k6 version: v%s", fullVersion())
	return nil
}

func Execute() {
	cobra.CheckErr(c.rootCmd.Execute())
}

func init() {
}

func (c *rootCommand) initLogger() error {

	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	if c.debug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("Debug mode enabled")
	}
}

func (c *rootCommand) initConfig() error {
	if c.debug {
		c.Logger.SetLevel(logrus.DebugLevel)
	}

	// if cfgFile != "" {
	// 	viper.SetConfigFile(cfgFile)
	// } else {
	// 	viper.AddConfigPath(".")
	// 	viper.SetConfigName(".sonarcli")
	// }

	// viper.AutomaticEnv()

	// if err := viper.ReadInConfig(); err == nil {
	// 	logrus.Info("Using config file:", viper.ConfigFileUsed())
	// }
}
