package main

import (
	"encoding/json"
	"os"
	. "sonaqube-analysis/internal/common"
	"sonaqube-analysis/internal/git"
	"sonaqube-analysis/internal/sonar"

	"github.com/sirupsen/logrus"
)

type Configuration struct {
	SonarURL       string   `json:"sonarUrl"`
	ProjectKey     string   `json:"projectKey"`
	AuthToken      string   `json:"authToken"`
	TwentyRepoURL  string   `json:"twentyRepoURL"`
	TwentyRepoPath string   `json:"twentyRepoPath"`
	Tags           []string `json:"tags"`
}

func main() {
	// Configurações do projeto
	// twentyRepoURL := "https://github.com/twentyhq/twenty.git"
	// twentyRepoPath := "./repo/twenty"
	// tags := []string{"v0.33.3", "v0.33.2", "v0.33.1", "v0.33.0"} // Exemplos de tags
	// sonarURL := "http://localhost:9000"
	// projectKey := "Test"
	// authToken := "sqa_9d87f7e834accb935cfafdc9a8881ebb4dc0e149" // Define o token como variável de ambiente

	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	file, errOpen := os.Open("./config/config.json")
	if errOpen != nil {
		logger.Error("error Open:", errOpen)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		logger.Error("error Decode:", err)
	}

	Info("Configurações do projeto:")
	logger.Infof("Repo URL: %s", configuration.TwentyRepoURL)
	logger.Infof("Repo Path: %s", configuration.TwentyRepoPath)
	logger.Infof("Tags: %v", configuration.Tags)
	logger.Infof("SonarQube URL: %s", configuration.SonarURL)
	logger.Infof("Project Key: %s", configuration.ProjectKey)
	logger.Infof("Auth Token: %s", configuration.AuthToken)

	// Clonar repositório
	if err := git.CloneRepo(configuration.TwentyRepoURL, configuration.TwentyRepoPath); err != nil {
		if err.Error() == "repository already exists" {
			logger.Warn("Repositório já existe, continuando...")
		} else {
			logger.Fatalf("Erro ao clonar repositório: %v", err)
		}
	}

	for _, tag := range configuration.Tags {
		logger.Infof("Analisando tag: %s", tag)

		// Alterar para a tag específica
		if err := git.CheckoutTag(configuration.TwentyRepoPath, tag); err != nil {
			Info("Erro ao alternar para a tag %s: %v", tag, err)
			// log.Printf("Erro ao alternar para a tag %s: %v", tag, err)
			continue
		}

		if err := sonar.CreateSonarProject("twentyTest", "TwentyTestProject"); err != nil {
			logger.Errorf("Erro ao criar um projeto no SonarQUbe: %v", err)
			continue
		}

		// Rodar SonarScanner
		// if err := sonar.RunSonarScanner(twentyRepoPath); err != nil {
		// 	log.Printf("Erro ao executar SonarScanner na tag %s: %v", tag, err)
		// 	continue
		// }

		// Obter métricas via API do SonarQube
		// metrics, err := sonar.GetSonarMetrics(sonarURL, projectKey, authToken)
		// if err != nil {
		// 	log.Printf("Erro ao obter métricas para a tag %s: %v", tag, err)
		// 	continue
		// }

		// Exportar métricas para CSV
		// outputFile := filepath.Join("./results", fmt.Sprintf("metrics_%s.csv", tag))

		// Converter métricas para o tipo correto
		// var dataMetrics []data.Metric
		// for _, m := range metrics {
		// 	dataMetrics = append(dataMetrics, data.Metric{
		// 		Name:    m.Key,
		// 		Measure: m.Measure,
		// 	})
		// }

		// if err := data.ExportToCSV(dataMetrics, outputFile); err != nil {
		// 	log.Printf("Erro ao exportar métricas da tag %s: %v", tag, err)
		// }
	}
}
