/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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

	client := &http.Client{Jar: AuthData.CookieJar}
	req, err := http.NewRequest("GET", "http://localhost:9000/api/projects/search", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("X-XSRF-TOKEN", AuthData.XSRFToken)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Scan failed with status: %s\n", resp.Status)
	} else {
		fmt.Printf("Scan successful with status: %s\n", resp.Status)
	}
	body, _ := io.ReadAll(resp.Body)
	// fmt.Printf("Response body: %s\n", body)

	type Component struct {
		Key string `json:"key"`
	}

	type Response struct {
		Components []Component `json:"components"`
	}

	var result Response
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}

	var keys []string
	for _, component := range result.Components {
		// fmt.Printf("Project Key: %s\n", component.Key)
		keys = append(keys, component.Key)
	}

	for _, key := range keys {
		deleteURL := fmt.Sprintf("http://localhost:9000/api/projects/delete?project=%s", key)
		req, err := http.NewRequest("POST", deleteURL, nil)
		if err != nil {
			fmt.Printf("Error creating delete request for %s: %v\n", key, err)
			continue
		}
		req.Header.Set("X-XSRF-TOKEN", AuthData.XSRFToken)

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error deleting project %s: %v\n", key, err)
			continue
		}
		resp.Body.Close()

		if resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusOK {
			fmt.Printf("Deleted project: %s\n", key)
		} else {
			fmt.Printf("Failed to delete project %s: %s\n", key, resp.Status)
		}
	}
}
