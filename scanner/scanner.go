package scanner

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"

	"github.com/fatih/color"
)

var (
	CompiledPatterns map[string]*regexp.Regexp
	Reports          []Report
	LeakColor        = color.New(color.FgHiMagenta).Add(color.Bold)
	ErroColor        = color.New(color.FgHiRed).Add(color.Bold)
	InfoColor        = color.New(color.FgHiBlue).Add(color.Bold)
	SucessColor      = color.New(color.FgHiGreen).Add(color.Bold)
	hasErrors        = false
)

func ReceiveRepo() string {
	repoPath := flag.String("repo", "", "Caminho para o Repositório")
	flag.Parse()

	if *repoPath == "" {
		ErroColor.Print("[ERRO] ")
		fmt.Println("Você deve fornecer o caminho para o repositório usando --repo")
		hasErrors = true
		os.Exit(1)
	}

	fmt.Printf("Caminho para o repositório: %s\n", *repoPath)
	return *repoPath
}

func ReadPath(repoPath string) error {
	InfoColor.Print("[INFO] ")
	fmt.Printf("Iniciando a leitura do repositório: %s\n", repoPath)
	err := filepath.WalkDir(repoPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			ErroColor.Print("[ERRO] ")
			fmt.Printf("Não foi possível acessar o caminho: %s\n", path)
			hasErrors = true
			return err
		}
		if d.IsDir() {
			if d.Name() == ".git" {
				return filepath.SkipDir
			}
		} else {
			if err := ReadFile(path, repoPath); err != nil {
				ErroColor.Print("[ERRO] ")
				fmt.Printf("Falha ao ler o arquivo: %s\n%v\n", path, err)
				hasErrors = true
				return err
			}
		}
		return nil
	})

	if err != nil {
		ErroColor.Print("[ERRO] ")
		fmt.Printf("Ocorreu um erro durante a varredura do Repositório: %s\n%v\n", repoPath, err)
		hasErrors = true
		return err
	}

	return nil
}

func ReadFile(archive string, repoPath string) error {
	file, err := os.Open(archive)
	if err != nil {
		ErroColor.Print("[ERRO] ")
		fmt.Printf("Não foi possível abrir o arquivo %s\n", archive)
		hasErrors = true
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
					ErroColor.Print("[ERRO] ")
					fmt.Printf("Não foi possível calcular o caminho relativo: %v\n", err)
					hasErrors = true
					continue
				}
				LeakColor.Print("[ATENÇÃO] ")
				fmt.Printf("Vazamento Encontrado em %s: %s\n Linha: %d\n --> %s\n\n", relativePath, name, line, lineContent)
				break
			}
		}
		line++
	}

	if err := scanner.Err(); err != nil {
		ErroColor.Print("[ERRO] ")
		fmt.Printf("Não foi possível ler o arquivo: %s\n%v", archive, err)
		hasErrors = true
		return err
	}

	return nil
}

func FinalizeReports(repoName string) error {

	if hasErrors {
		InfoColor.Print("[INFO] ")
		fmt.Println("Algum erro ocorreu durante o processo, portanto, o relatório não foi gerado")
		return nil
	}
	if len(Reports) == 0 {
		InfoColor.Print("[INFO] ")
		fmt.Println("Nenhum vazamento encontrado :)")
		InfoColor.Print("[INFO] ")
		fmt.Println("O Relatório não será gerado!")
		return nil
	}

	err := MakeReports(Reports, filepath.Base(repoName))
	if err != nil {
		ErroColor.Print("[ERRO] ")
		fmt.Printf("Não foi possível salvar o relatório: %v", err)
		return err
	}
	return nil
}
