/*
Copyright © 2025 Joao Cervino jcervinobarbosa@gmail.com
*/
package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var AuthData *AuthResponse

// AuthResponse contains extracted auth cookies
type AuthResponse struct {
	XSRFToken  string
	JWTSession string
	CookieJar  http.CookieJar
}

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
	authResp, err := auth()
	if err != nil {
		log.Fatalf("Authentication failed: %v", err)
	}

	healthCheck(authResp)
}

func auth() (*AuthResponse, error) {
	loginURL := "http://localhost:9000/api/authentication/login"

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %w", err)
	}

	client := &http.Client{
		Jar: jar,
	}

	form := url.Values{}
	form.Set("login", "admin")
	form.Set("password", "zy3fnVnvKLw4dca!")

	req, err := http.NewRequest("POST", loginURL, io.NopCloser(strings.NewReader(form.Encode())))
	if err != nil {
		return nil, fmt.Errorf("failed to create login request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36")
	req.Header.Set("Origin", "http://localhost:9001")
	req.Header.Set("Referer", "http://localhost:9001/projects/create?mode=manual&setncd=true")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("login request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("login failed: %s\n%s", resp.Status, string(body))
	}

	// Extract cookies from the response
	var xsrfToken, jwtSession string
	for _, cookie := range jar.Cookies(req.URL) {
		switch cookie.Name {
		case "XSRF-TOKEN":
			xsrfToken = cookie.Value
		case "JWT-SESSION":
			jwtSession = cookie.Value
		}
	}

	if xsrfToken == "" || jwtSession == "" {
		return nil, fmt.Errorf("required auth cookies not found")
	}

	return &AuthResponse{
		XSRFToken:  xsrfToken,
		JWTSession: jwtSession,
		CookieJar:  jar,
	}, nil
}

func healthCheck(authResp *AuthResponse) {
	client := &http.Client{Jar: authResp.CookieJar}

	req, err := http.NewRequest("GET", "http://localhost:9000/api/system/health", nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("User-Agent", "insomnia/10.3.1")
	req.Header.Set("X-XSRF-TOKEN", authResp.XSRFToken)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Health check request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Health check failed: %s\nResponse: %s", resp.Status, body)
	}

	fmt.Printf("Health check passed: %s\n", body)
}
