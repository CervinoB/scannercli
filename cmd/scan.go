/*
Copyright © 2025 Joao Cervino jcervinobarbosa@gmail.com
*/
package cmd

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

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
	fmt.Printf("Response body: %s\n", body)
	timeStamp := time.Now().Unix()

	if err := createNewProject(client, "new-project-"+fmt.Sprintf("%d", timeStamp)); err != nil {
		fmt.Printf("Error creating new project: %v\n", err)
		return
	}

	fmt.Println("Scan completed")
}

func createNewProject(client *http.Client, projectName string) error {
	endpoint := "http://localhost:9000/api/projects/create"

	form := url.Values{}
	form.Set("creationMode", "manual")
	form.Set("monorepo", "false")
	form.Set("project", projectName)
	form.Set("name", projectName)
	form.Set("mainBranch", "main")

	req, err := http.NewRequest("POST", endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-XSRF-TOKEN", AuthData.XSRFToken)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send project creation request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("project creation failed: status %d, response: %s", resp.StatusCode, body)
	}

	fmt.Printf("Project '%s' created successfully\n", projectName)
	return nil
}
