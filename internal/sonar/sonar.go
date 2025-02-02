package sonar

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
)

// Project represents a SonarQube project with its name and key.
type Project struct {
	Name    string `json:"name"`    // Name is the name of the project.
	Project string `json:"project"` // Project is the key of the project in SonarQube.
}

// RunSonarScanner runs the SonarQube scanner on the specified project path with the given project key.
// projectPath: the path to the project to be scanned.
// projectKey: the key of the project in SonarQube.
func RunSonarScanner(projectPath string, projectKey string) error {
	cmd := exec.Command("sonar-scanner",
		"-Dsonar.projectKey="+projectKey,
		"-Dsonar.sources="+projectPath,
		"-Dsonar.host.url=http://localhost:9000",
		"-Dsonar.token=sqa_9d87f7e834accb935cfafdc9a8881ebb4dc0e149")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// CreateSonarProject creates a new project in SonarQube with the given name and project key.
// Parameters:
// - name: The name of the project.
// - project: The project key.
func CreateSonarProject(name string, project string) error {
	url := fmt.Sprintf("http://localhost:9000/api/projects/create?name=%s&project=%s", name, project)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer sqa_9d87f7e834accb935cfafdc9a8881ebb4dc0e149")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to create project, status code: %d", resp.StatusCode)
	}

	return nil
}
