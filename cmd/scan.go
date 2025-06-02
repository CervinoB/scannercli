/*
Copyright © 2025 Joao Cervino jcervinobarbosa@gmail.com
*/
package cmd

import (
	"fmt"
	"io"
	"net/http"

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

	fmt.Println("Scan completed")
}
