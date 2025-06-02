/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/CervinoB/scannercli/internal/api"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: deleteRun,
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func deleteRun(cmd *cobra.Command, args []string) {
	fmt.Println("delete called")

	keys, err := api.ListProjects("http://localhost:9000", AuthData)
	if err != nil {
		fmt.Println("Failed to list projects:", err)
		return
	}

	for _, key := range keys {
		if err := api.DeleteProject("http://localhost:9000", key, AuthData); err != nil {
			fmt.Printf("Failed to delete %s: %v\n", key, err)
		} else {
			fmt.Printf("Deleted project: %s\n", key)
		}
	}
}
