package main

import (
	"LeakFinder/scanner"
	"fmt"
	"log"
)

func main() {
	repoPath := scanner.ReceiveRepo()

	var err error
	scanner.CompiledPatterns, err = scanner.GetPatterns()
	if err != nil {
		log.Fatalf("[ERRO] Falha ao compilar padrões: %v\n", err)
	}

	scanner.ReadPath(repoPath)
	err = scanner.FinalizeReports(repoPath)
	if err != nil {
		fmt.Printf("[ERRO] %v\n", err)
	}
}
