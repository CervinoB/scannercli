package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/CervinoB/scannercli/internal/logging"
	"github.com/spf13/viper"
)

func ListProjects(baseURL string, auth *AuthResponse) ([]string, error) {
	client := &http.Client{Jar: auth.CookieJar}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/projects/search", baseURL), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-XSRF-TOKEN", auth.XSRFToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %s", resp.Status)
	}

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		Components []struct {
			Key string `json:"key"`
		} `json:"components"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	var keys []string
	for _, c := range result.Components {
		keys = append(keys, c.Key)
	}
	return keys, nil
}

func CreateProject(baseURL, name string, auth *AuthResponse) error {
	if FF_CLEANUP := viper.GetString("FF_CLEANUP"); FF_CLEANUP == "true" {
		DeleteProject(baseURL, name, auth)
		logging.Logger.Infof("FF_CLEANUP enabled: deleted existing project before creating new one (key: %s)", name)
	}

	client := &http.Client{Jar: auth.CookieJar}
	endpoint := fmt.Sprintf("%s/api/projects/create", baseURL)

	form := url.Values{}
	form.Set("creationMode", "manual")
	form.Set("monorepo", "false")
	form.Set("project", name)
	form.Set("name", name)
	form.Set("mainBranch", "main")

	req, err := http.NewRequest("POST", endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-XSRF-TOKEN", auth.XSRFToken)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create project: %s\n%s", resp.Status, string(body))
	}

	return nil
}

func DeleteProject(baseURL, key string, auth *AuthResponse) error {
	client := &http.Client{Jar: auth.CookieJar}
	url := fmt.Sprintf("%s/api/projects/delete?project=%s", baseURL, key)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-XSRF-TOKEN", auth.XSRFToken)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete project %s: %s", key, resp.Status)
	}

	return nil
}

func DeleteAllProjects(baseURL string, auth *AuthResponse) error {
	keys, err := ListProjects(baseURL, auth)
	if err != nil {
		logging.Logger.Errorf("Failed to list projects: %v\n", err)
		return err
	}

	for _, key := range keys {
		if err := DeleteProject(baseURL, key, auth); err != nil {
			logging.Logger.Errorf("Failed to delete %s: %v\n", key, err)
		} else {
			logging.Logger.Infof("Deleted project: %s\n", key)
		}
	}
	return nil
}

func GenerateAnalysisToken(baseURL, key string, auth *AuthResponse) (string, error) {
	client := &http.Client{Jar: auth.CookieJar}
	endpoint := fmt.Sprintf("%s/api/user_tokens/generate", baseURL)

	form := url.Values{}
	form.Set("name", fmt.Sprint("Analysis Token for ", key))
	form.Set("type", "PROJECT_ANALYSIS_TOKEN")
	form.Set("projectKey", key)
	form.Set("expirationDate", "2026-06-24")

	req, err := http.NewRequest("POST", endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-XSRF-TOKEN", auth.XSRFToken)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to generate analysis token: %s\n%s", resp.Status, string(body))
	}

	var result struct {
		Login          string `json:"login"`
		Name           string `json:"name"`
		Token          string `json:"token"`
		CreatedAt      string `json:"createdAt"`
		Type           string `json:"type"`
		ProjectKey     string `json:"projectKey"`
		ExpirationDate string `json:"expirationDate"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Token, nil
}
