/*
Copyright © 2025 Joao Cervino jcervinobarbosa@gmail.com
*/
package cmd

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
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

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(showConfigCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
