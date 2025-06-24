/*
Copyright © 2025 Joao Cervino jcervinobarbosa@gmail.com
*/
package cmd

import (
	"os"
	"sort"

	"github.com/CervinoB/scannercli/internal/api"
	"github.com/CervinoB/scannercli/internal/export"
	"github.com/CervinoB/scannercli/internal/git"
	"github.com/CervinoB/scannercli/internal/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	logging.Logger.Println("scan called")

	name, url, sonarHost := getConfigValues()

	err := api.CreateProject(sonarHost, name, AuthData)
	if err != nil {
		logging.Logger.Errorf("Error creating project: %v\n", err)
		return
	} else {
		logging.Logger.Printf("Project created with key: %s\n", name)
	}

	err = git.CloneRepository(url, repoPath+"/"+name, Debug)
	if err != nil {
		logging.Logger.Printf("Error cloning repository: %v\n", err)
		return
	}

	token, err := api.GenerateAnalysisToken(sonarHost, name, AuthData)
	if err != nil {
		logging.Logger.Printf("Error generating analysis token: %v\n", err)
		return
	}

	logging.Logger.Printf("Analysis token generated: %s\n", token)

	tagList, err := getTags(name)
	if err != nil {
		logging.Logger.Printf("Error getting tags: %v\n", err)
		return
	}

	if len(tagList) > 0 {
		firstTag := tagList[0]
		err := git.CheckoutTag(repoPath+"/"+name, firstTag, Debug)
		if err != nil {
			logging.Logger.Printf("Error checking out tag %s: %v\n", firstTag, err)
			return
		}
		logging.Logger.Printf("Checked out tag: %s\n", firstTag)
	} else {
		logging.Logger.Println("No tags found to scan")
	}

	//TODO: Implement the scanner logic here using the token to run
	/*
			sonar-scanner \
		  -Dsonar.projectKey=twenty \
		  -Dsonar.sources=. \
		  -Dsonar.host.url=http://localhost:9000 \
		  -Dsonar.token=sqp_08d6bd9df2a4365c21a1c8af38c5bfaa0d416558
	*/

	for _, tag := range tagList {
		err := git.CheckoutTag(repoPath+"/"+name, tag, Debug)
		if err != nil {
			logging.Logger.Printf("Error checking out tag %s: %v\n", tag, err)
			return
		}
		//exec command
		err = api.ExecSonarScanner(name, token, sonarHost, repoPath+"/"+name, Debug)
		if err != nil {
			logging.Logger.Printf("Error running sonar-scanner: %v\n", err)
			return
		}
		logging.Logger.Printf("Sonar scan completed for tag: %s\n", tag)
		err = exportAllIssues(name, tag, sonarHost)
		if err != nil {
			logging.Logger.Printf("Error exporting issues for tag %s: %v\n", tag, err)
			// return
		}
	}

	logging.Logger.Info("Scan completed")
}

func getTags(name string) ([]string, error) {
	tagList, err := git.ListTags(repoPath + "/" + name)
	if err != nil {
		logging.Logger.Printf("Error listing tags: %v\n", err)
		return nil, err
	}
	logging.Logger.Debugf("Tags found: %v\n", tagList)

	if len(tagList) > 1 {
		// Sort tags alphabetically (or by semver if needed)
		// For now, sort alphabetically
		sortedTags := make([]string, len(tagList))
		copy(sortedTags, tagList)
		// You can use sort.Strings for alphabetical order
		// import "sort" at the top if not already imported
		sort.Strings(sortedTags)
		tagList = sortedTags
	}
	return tagList, nil
}

func getConfigValues() (string, string, string) {
	name := viper.GetString("name")
	url := viper.GetString("url")
	sonarHost := viper.GetString("sonarHost")
	return name, url, sonarHost
}

func exportAllIssues(name, version, sonarHost string) error {
	Issues, err := api.ReadAllIssues(name, sonarHost, AuthData)
	if err != nil {
		logging.Logger.Errorf("Error reading issues: %v\n", err)
		// return err
	}
	logging.Logger.Infof("Issues found: %d\n", len(Issues))

	// write to CSV file
	if len(Issues) == 0 {
		logging.Logger.Info("No issues found.")
		return nil
	}
	csvData, err := export.ExportCSV(Issues)
	if err != nil {
		logging.Logger.Errorf("Error exporting issues to CSV: %v\n", err)
		return err
	}
	logging.Logger.Infof("CSV data generated successfully:\n%s\n", csvData)

	// path := "./data/" + name + "/issues.csv"

	os.MkdirAll("./data/"+name, 0755) // Ensure the directory exists
	if err := os.WriteFile("./data/"+name+"/issues-"+version+".csv", []byte(csvData), 0644); err != nil {
		logging.Logger.Errorf("Error writing CSV file: %v\n", err)
		return err
	}
	logging.Logger.Infof("CSV file written successfully: ./data/%s/issues-%s.csv\n", name, version)
	return nil
}
