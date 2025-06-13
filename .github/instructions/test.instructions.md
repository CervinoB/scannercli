---
description: Generate an implementation plan for new features or refactoring existing code.
tools: ["codebase", "fetch", "findTestFiles", "githubRepo", "search", "usages"]
---

Create a Golang CLI application using Cobra and Viper that:

1. **Configuration**:

   - Reads a YAML config file (`config.yml`) using Viper with the following structure:

     ```yaml
     repositories:
       - url: "https://github.com/owner/repo1.git"
         path: "./temp/repo1"
       - url: "https://github.com/owner/repo2.git"
         path: "./temp/repo2"

     sonarqube:
       url: "http://sonarqube.example.com"
       token: "your_sonarqube_token"
       project_key_prefix: "temp_"

     output:
       csv_file: "sonarqube_metrics.csv"
     ```

2. **Functionality**:

   - For each repository in the config:
     a. Clones the repository if not already present
     b. Fetches all git tags and checks out the first tag (sorted by semantic versioning)
     c. Runs SonarQube scanner with generated project key (`prefix + repo_name`)
     d. Waits for scan to complete and fetches metrics:
     - Bugs
     - Vulnerabilities
     - Code Smells
     - Coverage
     - Duplicated Lines
     - Technical Debt
       e. Records metrics along with repository name and tag version

3. **Output**:

   - Saves all metrics to a CSV file with columns:
     - Repository
     - Tag
     - ScanDate
     - Bugs
     - Vulnerabilities
     - CodeSmells
     - Coverage
     - DuplicatedLines
     - TechnicalDebt

4. **Implementation Guidelines**:

   - Use Cobra for CLI structure with these commands:
     - `scan`: Runs the full scanning process
     - `list`: Lists configured repositories
     - `version`: Shows tool version
   - Keep it simple - no need for extensive error handling or logging
   - Use simple Git commands via `exec.Command`
   - For SonarQube interaction:
     - Use the scanner CLI for running scans
     - Use the REST API (via HTTP client) to fetch metrics
   - Assume SonarQube scanner is already installed on the system
   - No authentication for Git repos (all public)

5. **Example Usage**:

   ```bash
   # Run all scans
   ./sonarscanner scan --config config.yml

   # List configured repos
   ./sonarscanner list --config config.yml
   ```
