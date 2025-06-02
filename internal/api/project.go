package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
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
