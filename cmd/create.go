/*
Copyright © 2025 Joao Cervino jcervinobarbosa@gmail.com
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a default configuration file",
	Long: `Create a default YAML configuration file for SonarCLI.

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create called")
		config := `# sonarcli configuration
repositories:
  - url: "https://github.com/example/repo"
    tag: "v1.0.0"
scanners:
  - sonarqube
  - eslint
docker:
		enabled: true
`
		err := os.WriteFile("sonarcli-config.yml", []byte(config), 0644)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("Failed to create config file")
			return
		}

		logrus.Info("Configuration file created: sonarcli-config.yml")
	},
}

func init() {
	configCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
