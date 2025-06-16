package export

import (
	"encoding/csv"
	"strconv"
	"strings"

	"github.com/CervinoB/scannercli/internal/api"
)

func ExportCSV(issues []api.Issue) (string, error) {
	if len(issues) == 0 {
		return "", nil
	}

	var sb strings.Builder
	writer := csv.NewWriter(&sb)

	// Write header
	header := []string{"Component", "Line", "Severity", "Type", "Message", "Effort", "Author", "SoftwareQuality"}
	if err := writer.Write(header); err != nil {
		return "", err
	}

	// Write each issue
	for _, issue := range issues {
		softwareQuality := ""
		if len(issue.Impacts) > 0 {
			softwareQuality = issue.Impacts[0].SoftwareQuality
		}

		record := []string{
			issue.Component,
			strconv.Itoa(issue.Line),
			issue.Severity,
			issue.Type,
			issue.Message, // Automatic escaping handled here
			issue.Effort,
			issue.Author,
			softwareQuality,
		}

		if err := writer.Write(record); err != nil {
			return "", err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", err
	}

	return sb.String(), nil
}
