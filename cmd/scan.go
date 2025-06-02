/*
Copyright © 2025 Joao Cervino jcervinobarbosa@gmail.com
*/
package cmd

import (
	"fmt"
	"time"

	"github.com/CervinoB/scannercli/internal/api"
	"github.com/spf13/cobra"
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
	fmt.Println("scan called")

	timestamp := time.Now().Unix()
	projectName := fmt.Sprintf("new-project-%d", timestamp)

	if err := api.CreateProject("http://localhost:9000", projectName, AuthData); err != nil {
		fmt.Printf("Error creating project: %v\n", err)
		return
	}
	fmt.Println("Scan completed")
}
