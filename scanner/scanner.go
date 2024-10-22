package scanner

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
)

var CompiledPatterns map[string]*regexp.Regexp
var Reports []Report

func ReceiveRepo() string {
	repoPath := flag.String("repo", "", "Caminho para o Repositório")
	flag.Parse()

	if *repoPath == "" {
		fmt.Println("[ERRO] Você deve fornecer o caminho para o repositório usando --repo")
		os.Exit(1)
	}

	fmt.Printf("Caminho para o repositório: %s\n", *repoPath)
	return *repoPath
}

func ReadPath(repoPath string) {
	fmt.Printf("[INFO] Iniciando a leitura do repositório: %s\n", repoPath)
	err := filepath.WalkDir(repoPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("[ERRO] Não foi possível acessar o caminho: %s\n", path)
			return err
		}
		if d.IsDir() {
			if d.Name() == ".git" {
				return filepath.SkipDir
			}
		} else {
			if err := ReadFile(path, repoPath); err != nil {
				fmt.Printf("[ERRO] Falha ao ler o arquivo: %s\n%v\n", path, err)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("[ERRO] Ocorreu um erro durante a varredura do Repositório: %s\n%v\n", repoPath, err)
	}
}

func ReadFile(archive string, repoPath string) error {
	file, err := os.Open(archive)
	if err != nil {
		fmt.Printf("[ERRO] Não foi possível abrir o arquivo %s\n", archive)
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	const maxCapacity = 1024 * 1024
	scanner.Buffer(make([]byte, 0, maxCapacity), maxCapacity)
	line := 1

	for scanner.Scan() {
		lineContent := scanner.Text()

		for name, pattern := range CompiledPatterns {
			if pattern.MatchString(lineContent) {
				relativePath, err := filepath.Rel(repoPath, archive)
				Reports = append(Reports, Report{
					FilePath: relativePath,
					LeakType: name,
					Line:     line,
					Content:  lineContent,
				})
				if err != nil {
					fmt.Printf("[ERRO] Não foi possível calcular o caminho relativo: %v\n", err)
					continue
				}
				fmt.Printf("[ATENÇÃO] Vazamento Encontrado em %s: %s\n Linha: %d\n --> %s\n\n", relativePath, name, line, lineContent)
				break
			}
		}
		line++
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("[ERRO] Não foi possível ler o arquivo: %s\n%v", archive, err)
	}

	return nil
}

func FinalizeReports(repoName string) error {
	err := MakeReports(Reports, filepath.Base(repoName))
	if err != nil {
		return fmt.Errorf("[ERRO] Não foi possível salvar o relatório: %v", err)
	}
	return nil
}
