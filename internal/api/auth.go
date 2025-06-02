package api

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

// AuthResponse contains extracted auth cookies and session
type AuthResponse struct {
	XSRFToken  string
	JWTSession string
	CookieJar  http.CookieJar
}

// Authenticate logs into the system and returns AuthResponse
func Authenticate(loginURL, username, password string) (*AuthResponse, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %w", err)
	}

	client := &http.Client{Jar: jar}

	form := url.Values{}
	form.Set("login", username)
	form.Set("password", password)

	req, err := http.NewRequest("POST", loginURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create login request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("login request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("login failed: %s\n%s", resp.Status, string(body))
	}

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
