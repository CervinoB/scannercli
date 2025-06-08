/*
Copyright © 2025 Joao Cervino jcervinobarbosa@gmail.com
*/
package cmd

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/CervinoB/scannercli/internal/api"
	"github.com/CervinoB/scannercli/internal/git"
	"github.com/CervinoB/scannercli/internal/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan repositories with the scanner",
	Long: `The scan command allows you to scan repositories using the scanner tool.
It provides detailed analysis and insights for the specified repositories.`,
	Run: scanRun,
}

func init() {
	rootCmd.AddCommand(scanCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func scanRun(cmd *cobra.Command, args []string) {
	logging.Logger.Println("scan called")

	for key, value := range viper.GetViper().AllSettings() {
		logging.Logger.WithFields(log.Fields{
			key: value,
		}).Info("Command Flag")
	}

	timestamp := time.Now().Unix()
	projectName := fmt.Sprintf("new-project-%d", timestamp)

	err := api.CreateProject("http://localhost:9000", projectName, AuthData)
	if err != nil {
		logging.Logger.Errorf("Error creating project: %v\n", err)
		return
	} else {
		logging.Logger.Printf("Project created with key: %s\n", projectName)
	}

	err = git.CloneRepository("https://github.com/twentyhq/twenty.git", repoPath+"/"+projectName)
	if err != nil {
		logging.Logger.Printf("Error cloning repository: %v\n", err)
		return
	}

	logging.Logger.Info("Scan completed")
}
