/*
Copyright © 2025 Joao Cervino jcervinobarbosa@gmail.com
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/CervinoB/scannercli/internal/api"
	"github.com/CervinoB/scannercli/internal/export"
	"github.com/CervinoB/scannercli/internal/logging"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:        "list",
	Short:      "List the todos",
	Long:       "Listing the todos",
	Deprecated: "This command is deprecated and will be removed in future versions. Used only to test and setup project.",

	Run: listRun,
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func listRun(cmd *cobra.Command, args []string) {
	fmt.Println("list called")

	name, _, sonarHost := getConfigValues()
	version := "v0.1.0"

	// Issues, err := api.ReadIssues(name, sonarHost, AuthData)
	Issues, err := api.ReadAllIssues(name, sonarHost, AuthData)
	if err != nil {
		logging.Logger.Errorf("Error reading issues: %v\n", err)
		// return
	}
	logging.Logger.Infof("Issues found: %d\n", len(Issues))

	// write to CSV file
	if len(Issues) == 0 {
		logging.Logger.Info("No issues found.")
		return
	}
	csvData, err := export.ExportCSV(Issues)
	if err != nil {
		logging.Logger.Errorf("Error exporting issues to CSV: %v\n", err)
		return
	}
	logging.Logger.Infof("CSV data generated successfully:\n%s\n", csvData)

	// path := "./data/" + name + "/issues.csv"

	os.MkdirAll("./data/"+name, 0755) // Ensure the directory exists
	if err := os.WriteFile("./data/"+name+"/issues-"+version+".csv", []byte(csvData), 0644); err != nil {
		logging.Logger.Errorf("Error writing CSV file: %v\n", err)
		return
	}
	logging.Logger.Infof("CSV file written successfully: ./data/%s/issues-%s.csv\n", name, version)
}
