package api

import (
	"fmt"
	"io"
	"net/http"
)

// CheckHealth pings the health endpoint
func CheckHealth(url string, auth *AuthResponse) error {
	client := &http.Client{Jar: auth.CookieJar}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	req.Header.Set("X-XSRF-TOKEN", auth.XSRFToken)
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("health check request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check failed: %s\n%s", resp.Status, string(body))
	}

	return nil
}
