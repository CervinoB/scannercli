/*
Copyright © 2025 Joao Cervino jcervinobarbosa@gmail.com
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var dataFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "scannercli",
	Short: "A CLI tool to scan repositories and gather metrics",
	Long: `ScannerCLI is a command-line tool designed to scan repositories and
retrieve metrics for each commit, tag, and hash using the desired scanner.

This tool helps developers and teams analyze repository history and extract
valuable insights efficiently.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	cwd, err := os.Getwd()
	if err != nil {
		log.Println("Unable to detect current directory. Please set data file using --datafile.")
	}

	// Default to ./.tridos.json in the current directory
	defaultFile := filepath.Join(cwd, ".tridos.json")
	rootCmd.PersistentFlags().StringVar(&dataFile, "datafile", defaultFile, "data file to store todos")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	fmt.Println("Using data file:", dataFile)
}
