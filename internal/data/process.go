// package data

// // Metric representa uma métrica coletada do SonarQube.
// type Metric struct {
// 	MetricKey string
// 	Value     string
// }

// // ProcessMetrics normaliza os dados coletados, converte valores para tipos numéricos e valida inconsistências.
// func ProcessMetrics(metrics []Metric) (map[string]float64, error) {
// 	if len(metrics) == 0 {
// 		return nil, errors.New("nenhuma métrica encontrada para processar")
// 	}

// 	processedData := make(map[string]float64)

// 	for _, metric := range metrics {
// 		// Tenta converter o valor para um número de ponto flutuante.
// 		value, err := strconv.ParseFloat(metric.Value, 64)
// 		if err != nil {
// 			return nil, fmt.Errorf("erro ao processar métrica '%s': valor inválido '%s'", metric.MetricKey, metric.Value)
// 		}

// 		// Adiciona a métrica normalizada ao mapa.
// 		processedData[metric.MetricKey] = value
// 	}

// 	return processedData, nil
// }

// // CompareMetrics compara as métricas entre duas versões e retorna as diferenças.
// func CompareMetrics(oldMetrics, newMetrics map[string]float64) map[string]float64 {
// 	differences := make(map[string]float64)

// 	for key, newValue := range newMetrics {
// 		oldValue, exists := oldMetrics[key]
// 		if !exists {
// 			// Métrica nova, diferença igual ao valor atual.
// 			differences[key] = newValue
// 		} else {
// 			// Calcula a diferença entre a métrica antiga e a nova.
// 			differences[key] = newValue - oldValue
// 		}
// 	}

// 	return differences
// }

// // AggregateMetrics realiza agregações simples (como média ou soma) nas métricas coletadas.
// func AggregateMetrics(metrics map[string]float64) map[string]float64 {
// 	aggregatedData := make(map[string]float64)

// 	// Exemplo: Adicionando uma métrica de exemplo para a soma de "bugs" e "code_smells".
// 	if bugs, exists := metrics["bugs"]; exists {
// 		if smells, exists := metrics["code_smells"]; exists {
// 			aggregatedData["bugs_and_smells_total"] = bugs + smells
// 		}
// 	}

// 	return aggregatedData
// }
