package scanner

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Report struct {
	FilePath string `json:"FilePath"`
	LeakType string `json:"LeakType"`
	Line     int    `json:"Line"`
	Content  string `json:"Content"`
}

// Função principal que orquestra o processo de geração de relatórios
func MakeReports(reports []Report, repoName string) error {
	groupedReports := groupReportsByFile(reports)
	reportPath, file, err := createReportFile(repoName)
	if err != nil {
		return err
	}
	defer file.Close()

	return encodeReportsToJSON(groupedReports, file, reportPath)
}

func groupReportsByFile(reports []Report) map[string][]Report {
	groupedReports := make(map[string][]Report)
	for _, report := range reports {
		groupedReports[report.FilePath] = append(groupedReports[report.FilePath], report)
	}
	return groupedReports
}

func createReportFile(repoName string) (string, *os.File, error) {
	filename := fmt.Sprintf("%s_reports.json", repoName)
	reportPath := filepath.Join("reports", filename)

	file, err := os.OpenFile(reportPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		ErroColor.Print("[ERRO] ")
		fmt.Printf("Não foi possível criar ou abrir o arquivo de relatório: %v\n", err)
		return "", nil, err
	}

	return reportPath, file, nil
}

func encodeReportsToJSON(groupedReports map[string][]Report, file *os.File, reportPath string) error {
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")

	finalReports := prepareFinalReports(groupedReports)

	if err := encoder.Encode(finalReports); err != nil {
		ErroColor.Print("[ERRO] ")
		fmt.Printf("Não foi possível escrever no arquivo de relatório: %v\n", err)
		return err
	}

	SucessColor.Print("[SUCESS] ")
	fmt.Printf("Relatório salvo com sucesso em: %s\n", reportPath)
	return nil
}

func prepareFinalReports(groupedReports map[string][]Report) []struct {
	FilePath string `json:"File"`
	Leaks    []struct {
		LeakType string `json:"LeakType"`
		Line     int    `json:"Line"`
		Content  string `json:"Content"`
	} `json:"Leaks"`
} {
	var finalReports []struct {
		FilePath string `json:"File"`
		Leaks    []struct {
			LeakType string `json:"LeakType"`
			Line     int    `json:"Line"`
			Content  string `json:"Content"`
		} `json:"Leaks"`
	}

	leakSeen := make(map[string]struct{})

	for filePath, leaks := range groupedReports {
		var leakReports []struct {
			LeakType string `json:"LeakType"`
			Line     int    `json:"Line"`
			Content  string `json:"Content"`
		}

		for _, leak := range leaks {
			leakIdentifier := fmt.Sprintf("%s:%d:%s", leak.LeakType, leak.Line, leak.Content)

			if _, exists := leakSeen[leakIdentifier]; !exists {
				leakSeen[leakIdentifier] = struct{}{}
				leakReports = append(leakReports, struct {
					LeakType string `json:"LeakType"`
					Line     int    `json:"Line"`
					Content  string `json:"Content"`
				}{
					LeakType: leak.LeakType,
					Line:     leak.Line,
					Content:  leak.Content,
				})
			}
		}

		if len(leakReports) > 0 {
			finalReports = append(finalReports, struct {
				FilePath string `json:"File"`
				Leaks    []struct {
					LeakType string `json:"LeakType"`
					Line     int    `json:"Line"`
					Content  string `json:"Content"`
				} `json:"Leaks"`
			}{
				FilePath: filePath,
				Leaks:    leakReports,
			})
		}
	}

	return finalReports
}
