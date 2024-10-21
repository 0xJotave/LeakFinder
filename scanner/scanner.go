package scanner

import (
	"LeakFinder/config"
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

type ScanResult struct {
	FilePath string
	Line     int
	Content  string
}

func ScanRepo(repoPath string, config config.Config) []ScanResult {
	var results []ScanResult

	err := filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		for _, dir := range config.IgnoreDirs {
			if info.IsDir() && info.Name() == dir {
				return filepath.SkipDir
			}
		}

		if !info.IsDir() {
			scanFile(path, &results, config.Patterns)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error scanning repository: %v\n", err)
	}

	return results
}

func scanFile(filePath string, results *[]ScanResult, patterns map[string]string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", filePath, err)
		return
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		lineContent := scanner.Text()

		for name, pattern := range patterns {
			matched, err := regexp.MatchString(pattern, lineContent)
			if err != nil {
				fmt.Printf("Error compiling regex for %s: %v\n", name, err)
				continue
			}
			if matched {
				result := ScanResult{
					FilePath: filePath,
					Line:     lineNumber,
					Content:  lineContent,
				}
				*results = append(*results, result)
				break
			}
		}
	}

}
