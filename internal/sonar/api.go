// TODO: REVIEW
package sonar

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Metric struct {
	MetricKey string `json:"metric"`
	Value     string `json:"value"`
}

func GetSonarMetrics(baseURL, projectKey, authToken string) ([]Metric, error) {
	url := fmt.Sprintf("%s/api/measures/component?component=%s&metricKeys=bugs,code_smells", baseURL, projectKey)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Basic "+authToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("falha ao obter métricas: %s", resp.Status)
	}

	var result struct {
		Component struct {
			Measures []Metric `json:"measures"`
		} `json:"component"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	return result.Component.Measures, err
}
