package main

import (
	"LeakFinder/scanner"
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
}
