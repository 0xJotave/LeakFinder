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

func MakeReports(reports []Report, repoName string) error {
	groupedReports := make(map[string][]Report)

	for _, report := range reports {
		groupedReports[report.FilePath] = append(groupedReports[report.FilePath], report)
	}

	filename := fmt.Sprintf("%s_reports.json", repoName)
	reportPath := filepath.Join("reports", filename)

	file, err := os.OpenFile(reportPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("[ERRO] Não foi possível criar ou abrir o arquivo de relatório: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")

	var finalReports []struct {
		FilePath string `json:"File"`
		Leaks    []struct {
			LeakType string `json:"LeakType"`
			Line     int    `json:"Line"`
			Content  string `json:"Content"`
		} `json:"Leaks"`
	}

	// Para verificar se um vazamento já foi visto
	leakSeen := make(map[string]struct{})

	for filePath, leaks := range groupedReports {
		var leakReports []struct {
			LeakType string `json:"LeakType"`
			Line     int    `json:"Line"`
			Content  string `json:"Content"`
		}

		for _, leak := range leaks {
			// Criar um identificador único para cada vazamento
			leakIdentifier := fmt.Sprintf("%s:%d:%s", leak.LeakType, leak.Line, leak.Content)

			// Verificar se o vazamento já foi visto
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

		// Adicionar os vazamentos do arquivo processado
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

	if err = encoder.Encode(finalReports); err != nil {
		return fmt.Errorf("[ERRO] Não foi possível escrever no arquivo de relatório: %v", err)
	}

	fmt.Printf("Relatório salvo com sucesso em: %s\n", reportPath)
	return nil
}
