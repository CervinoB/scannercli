// // TODO: review
// package data

// import (
// 	"encoding/csv"
// 	"os"
// )

// type Metric struct {
// 	MetricKey string
// 	Value     string
// }

// func ExportToCSV(metrics []Metric, outputPath string) error {
// 	file, err := os.Create(outputPath)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	writer := csv.NewWriter(file)
// 	defer writer.Flush()

// 	// Cabeçalhos
// 	writer.Write([]string{"Metric", "Value"})

// 	// Dados
// 	for _, metric := range metrics {
// 		writer.Write([]string{metric.MetricKey, metric.Value})
// 	}

// 	return nil
// }
