/*
Copyright © 2025 Joao Cervino jcervinobarbosa@gmail.com
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/CervinoB/sonarcli/cmd/state"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type createCmd struct {
	gs *state.GlobalState
}

func (c *createCmd) run(cmd *cobra.Command, _ []string) error {
	var l logrus.FieldLogger = c.gs.Logger

	// defer func() {
	// 	if err == nil {
	// 		l.Debug("Everything has finished, exiting sonarcli normally!")
	// 	} else {
	// 		l.WithError(err).Debug("Everything has finished, exiting sonarcli with an error!")
	// 	}
	// }()
	// printBanner(c.gs)

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

	l.Info("Scan completed successfully.")
	return nil
}

func (c *createCmd) getCmdScan(gs *state.GlobalState) *cobra.Command {
	createCmd := &createCmd{gs: gs}
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a default configuration file",
		Long: `Create a default YAML configuration file for SonarCLI.

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		RunE: createCmd.run,
		// func(cmd *cobra.Command, args []string) {
		// 			fmt.Println("create called")
		// 			config := `# sonarcli configuration
		// repositories:
		//   - url: "https://github.com/example/repo"
		//     tag: "v1.0.0"
		// scanners:
		//   - sonarqube
		//   - eslint
		// docker:
		// 		enabled: true
		// `
		// 			err := os.WriteFile("sonarcli-config.yml", []byte(config), 0644)
		// 			if err != nil {
		// 				logrus.WithFields(logrus.Fields{
		// 					"error": err,
		// 				}).Error("Failed to create config file")
		// 				return
		// 			}

		// 			logrus.Info("Configuration file created: sonarcli-config.yml")
		// 		},
	}

	cmd.Flags().StringP("config", "f", "", "Configuration file to use")
	cmd.MarkFlagRequired("config")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	return cmd
}

// createCmd represents the create command
var createCmdOld = &cobra.Command{
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

// func init() {
// 	configCmd.AddCommand(createCmd)

// 	// Here you will define your flags and configuration settings.

// 	// Cobra supports Persistent Flags which will work for this command
// 	// and all subcommands, e.g.:
// 	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

// 	// Cobra supports local flags which will only run when this command
// 	// is called directly, e.g.:
// 	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
// }
