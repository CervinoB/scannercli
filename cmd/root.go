/*
Copyright © 2025 Joao Cervino jcervinobarbosa@gmail.com
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/CervinoB/scannercli/internal/api"
	"github.com/CervinoB/scannercli/internal/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var AuthData *api.AuthResponse

var repoPath string
var dataFile string
var Verbose bool
var Debug bool

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
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		sonarHost := viper.GetString("sonarHost")

		authResp, err := api.Authenticate(fmt.Sprintf("%s/api/authentication/login", sonarHost), "admin", "zy3fnVnvKLw4dca!")
		if err != nil {
			logging.Logger.Errorf("Authentication failed: %v", err)
			return fmt.Errorf("auth failed: %w", err)
		}
		AuthData = authResp
		if err := api.CheckHealth(fmt.Sprintf("%s/api/system/health", sonarHost), authResp); err != nil {
			logging.Logger.Errorf("Health check failed: %v", err)
			return fmt.Errorf("health check failed: %w", err)
		}

		return nil
	},
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

	cobra.OnInitialize(initConfig)

	cwd, err := os.Getwd()
	if err != nil {
		log.Println("Unable to detect current directory. Please set data file using --repoPath.")
	}

	// Default to ./repo/ in the current directory
	defaultPath := filepath.Join(cwd, "repo/")
	rootCmd.PersistentFlags().StringVar(&repoPath, "repoPath", defaultPath, "repository path")
	viper.BindPFlag("repoPath", rootCmd.PersistentFlags().Lookup("repoPath"))

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Display more verbose output in console output. (default: false)")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))

	rootCmd.PersistentFlags().BoolVarP(&Debug, "debug", "d", false, "Display debugging output in the console. (default: false)")
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	fmt.Println("Using repository path:", repoPath)
}

func initConfig() {
	logging.ConfigureLogger(Verbose, Debug)
	logging.Logger.Debug("Logger initialized")

	// Set the config file name and paths to search
	viper.SetConfigName("scannercli") // name of config file (without extension)
	viper.SetConfigType("yaml")       // or "json", "toml", etc.

	// Add paths to search for config files
	viper.AddConfigPath(".") // current directory

	// Read environment variables
	viper.AutomaticEnv()

	// If a config file is found, read it in
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	for key, value := range viper.GetViper().AllSettings() {
		logging.Logger.WithFields(map[string]interface{}{
			key: value,
		}).Info("Command Flag")
	}
}
