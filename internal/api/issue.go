package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Impact struct {
	SoftwareQuality string `json:"softwareQuality"`
	Severity        string `json:"severity"`
}

type Issue struct {
	Component string   `json:"component"`
	Line      int      `json:"line"`
	Severity  string   `json:"severity"`
	Type      string   `json:"type"`
	Message   string   `json:"message"`
	Effort    string   `json:"effort"`
	Author    string   `json:"author"`
	Impacts   []Impact `json:"impacts"`
}

func (i *Issue) PrettyPrint() string {
	return fmt.Sprintf(
		"Components: %s, Line: %d, Severity: %s, Type: %s, Message: %s, Effort: %s, Author: %s, Software Quality: %s",
		i.Component, i.Line, i.Severity, i.Type, i.Message, i.Effort, i.Author, i.Impacts[0].SoftwareQuality,
	)
}

func ReadIssues(projectKey, baseUrl string, auth *AuthResponse) ([]Issue, error) {
	// This function should read issues from a data source and return a slice of Issue structs.
	// For now, we will return an empty slice.
	client := &http.Client{Jar: auth.CookieJar}
	parameters := url.Values{}
	parameters.Add("components", projectKey)
	parameters.Add("s", "FILE_LINE")
	parameters.Add("issueStatuses", "CONFIRMED,OPEN")
	parameters.Add("additionalFields", "_all")
	parameters.Add("ps", "500")

	req, err := http.NewRequest("GET", baseUrl+"/api/issues/search?"+parameters.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-XSRF-TOKEN", auth.XSRFToken)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %s", resp.Status)
	}
	var result struct {
		Issues []Issue `json:"issues"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result.Issues, nil
}

func ReadAllIssues(projectKey, baseUrl string, auth *AuthResponse) ([]Issue, error) {
	// This function reads all issues for a project, handling pagination.
	var allIssues []Issue
	pageSize := 500
	page := 1

	for {
		issues, err := ReadIssuesWithPagination(projectKey, baseUrl, auth, page, pageSize)
		if err != nil {
			return allIssues, err
		}
		if len(issues) == 0 {
			break // No more issues to read
		}
		allIssues = append(allIssues, issues...)
		page++
	}

	return allIssues, nil
}

// ReadIssuesWithPagination reads issues for a project with pagination support.
func ReadIssuesWithPagination(projectKey, baseUrl string, auth *AuthResponse, page, pageSize int) ([]Issue, error) {
	client := &http.Client{Jar: auth.CookieJar}
	parameters := url.Values{}
	parameters.Add("components", projectKey)
	parameters.Add("s", "FILE_LINE")
	parameters.Add("issueStatuses", "CONFIRMED,OPEN")
	parameters.Add("additionalFields", "_all")
	parameters.Add("ps", fmt.Sprintf("%d", pageSize))
	parameters.Add("p", fmt.Sprintf("%d", page))

	req, err := http.NewRequest("GET", baseUrl+"/api/issues/search?"+parameters.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-XSRF-TOKEN", auth.XSRFToken)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %s", resp.Status)
	}
	var result struct {
		Issues []Issue `json:"issues"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result.Issues, nil
}
