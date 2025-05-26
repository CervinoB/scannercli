/*
Copyright © 2025 Joao Cervino jcervinobarbosa@gmail.com
*/
package cmd

import (
	"github.com/CervinoB/scannercli/cmd/state"
	"github.com/CervinoB/scannercli/internal/git"
	"github.com/CervinoB/scannercli/lib/consts"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type scanCmd struct {
	gs        *state.GlobalState
	scanner   string
	repoURL   string
	repoPath  string // path to the local repository
	clonePath string // path to clone the repository
	// docker bool
	tags []string // tags to apply to the scan
	name string   // name of the repository
}

func getCmdScan(gs *state.GlobalState) *cobra.Command {
	s := &scanCmd{gs: gs, clonePath: consts.ClonePath}
	scanCmd := &cobra.Command{
		Use:   "scan",
		Short: "Scan a repository for code smells",
		Long: `Scan a repository using SonarQube, ESLint, and other tools to detect code smells.

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		// Args: exactArgsWithMsg(1, "Repository URL to scan"),
		RunE: s.run,
	}

	scanCmd.Flags().StringP("repo", "t", "", "Repository URL to scan")
	scanCmd.Flags().StringP("commit", "c", "HEAD", "Commit hash or tag to analyze")
	scanCmd.Flags().StringP("scanner", "s", "", "Scanner to use (sonarqube, eslint, etc.)")

	return scanCmd
}

func (c *scanCmd) run(cmd *cobra.Command, args []string) (err error) {
	c.initConfig()

	if err = git.CloneRepo(c.gs, c.repoURL, c.clonePath+"/"+c.name); err != nil {
		if err.Error() == "repository already exists" {
			c.gs.Logger.Warn(err)
			c.gs.Logger.Debug("Try pulling...")
			git.Pull(c.gs, c.clonePath+"/"+c.name)
		} else {
			c.gs.Logger.Fatalf("Erro ao clonar repositório: %v", err)
		}
	}

	c.gs.Logger.Info("Starting repository scan...")

	// TODO: Implementar lógica de scan

	// aways run in dockerized modes

	// 2. Executar scanners (SonarQube, ESLint)

	// 3. Gerar relatório

	c.gs.Logger.Info("Scan completed successfully.")
	return nil
}

func (c *scanCmd) initConfig() {
	c.scanner = viper.GetString("scanner")
	c.gs.Logger.Infof("Using scanner: %s", c.scanner)

	c.tags = viper.GetStringSlice("tags")
	c.gs.Logger.Infof("Using tags: %s", c.tags)

	c.repoURL = viper.GetString("url")
	c.gs.Logger.Infof("Using URL: %s", c.repoURL)

	c.name = viper.GetString("name")
	c.gs.Logger.Infof("Using name: %s", c.name)

	c.repoPath = viper.GetString("nestPath")
	c.gs.Logger.Infof("Using nest server path: %s", c.repoPath)
}

func runDockerizedScan(target string) {
	logrus.Info("Running dockerized scanners", target)
	// TODO: Implementar lógica de Docker
}

func runLocalScan(target string) {
	logrus.Info("Running local scanners", target)
	// TODO: Chamar módulos de análise
}
