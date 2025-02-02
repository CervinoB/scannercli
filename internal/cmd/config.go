/*
Copyright © 2025 Joao Cervino jcervinobarbosa@gmail.com
*/
package cmd

import (
	"fmt"

	"github.com/CervinoB/sonarcli/cmd/state"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type configCmd struct {
	gs *state.GlobalState
}

func (c *configCmd) run(cmd *cobra.Command, _ []string) error {
	var l logrus.FieldLogger = c.gs.Logger

	// defer func() {
	// 	if err == nil {
	// 		l.Debug("Everything has finished, exiting sonarcli normally!")
	// 	} else {
	// 		l.WithError(err).Debug("Everything has finished, exiting sonarcli with an error!")
	// 	}
	// }()
	printBanner(c.gs)
	printBar(c.gs, nil)

	// TODO: Implementar lógica de config
	// 1. Carregar configurações
	target := "https://github.com/CervinoB/sonarcli"
	l.Infof("Starting scan for: %s", target)
	// if docker {
	// 	runDockerizedScan(target)
	// } else {
	// 	runLocalScan(target)
	// }

	// 2. Executar scanners (SonarQube, ESLint)

	// 3. Gerar relatório

	logrus.Info("Scan completed successfully.")
	return nil
}

func (c *configCmd) getCmdConfig(cmd *cobra.Command, _ []string) error {
	if c.gs.CfgFile != "" {
		viper.SetConfigFile(c.gs.CfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName(".sonarcli")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		logrus.Info("Using config file:", viper.ConfigFileUsed())
		return err
	}
	return nil

}

// configCmd represents the config command
var configCmdOld = &cobra.Command{
	Use:   "config",
	Short: "Manage SonarCLI configurations",
	Long:  `Manage configuration settings for SonarCLI.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("config called")
	},
}

var showConfigCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("Current configuration:")
		for _, key := range viper.AllKeys() {
			logrus.Infof("%s: %v", key, viper.Get(key))
		}
	},
}
