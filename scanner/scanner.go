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
	LeakColor   = color.New(color.FgHiMagenta).Add(color.Bold)
	ErroColor   = color.New(color.FgHiRed).Add(color.Bold)
	InfoColor   = color.New(color.FgHiBlue).Add(color.Bold)
	SucessColor = color.New(color.FgHiGreen).Add(color.Bold)
)

var (
	CompiledPatterns map[string]*regexp.Regexp
	Reports          []Report
	hasErrors        = false
)

func ReceiveRepo() string {
	repoPath := flag.String("repo", "", "Caminho para o Repositório")
	flag.Parse()

	if *repoPath == "" {
		HandleError("Você deve fornecer o caminho para o repositório usando --repo\n", "")
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
			HandleError("Não foi possível acessar o caminho: %s\n", path)
			return err
		}
		if d.IsDir() {
			if d.Name() == ".git" {
				return filepath.SkipDir
			}
		} else {
			if err := ReadFile(path, repoPath); err != nil {
				HandleError("Falha ao ler o arquivo: %s\n", path)
				return err
			}
		}
		return nil
	})
	if err != nil {
		HandleError("Ocorreu um erro durante a varredura do Repositório: %s\n", repoPath)
		return err
	}
	return nil
}

func ReadFile(archive string, repoPath string) error {
	file, err := os.Open(archive)
	if err != nil {
		HandleError("Não foi possível abrir o arquivo: %s\n", archive)
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	const maxCapacity = 1024 * 1024
	scanner.Buffer(make([]byte, 0, maxCapacity), maxCapacity)
	line := 1

	for scanner.Scan() {
		lineContent := scanner.Text()
		checkForLeaks(line, lineContent, repoPath, archive)
		line++

	}

	if err := scanner.Err(); err != nil {
		HandleError("Não foi possível ler o arquivo: %s\n", archive)
		return err
	}

	return nil
}

func checkForLeaks(line int, lineContent, repoPath, archive string) {
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
				HandleError("Não foi possível calcular o caminho relativo\n", "")
				continue
			}
			LeakColor.Print("[ATENÇÃO] ")
			fmt.Printf("Vazamento Encontrado em %s: %s\n Linha: %d\n --> %s\n\n", relativePath, name, line, lineContent)
			break
		}
	}
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
		HandleError("Não foi possível salvar o relatório\n", "")
		return err
	}
	return nil
}

func HandleError(message string, conteudo string) {
	ErroColor.Print("[ERRO] ")
	fmt.Printf(message, conteudo)
	hasErrors = true
}
