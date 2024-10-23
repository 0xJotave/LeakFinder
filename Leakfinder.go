package main

import (
	"LeakFinder/scanner"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/fatih/color"
)

var cRed = color.New(color.FgHiRed).Add(color.Bold)

func AsciiArt(wg *sync.WaitGroup) {
	defer wg.Done()

	asciiArt := `
 ██▓    ▓█████ ▄▄▄       ██ ▄█▀  █████▒██▓ ███▄    █ ▓█████▄ ▓█████  ██▀███  
▓██▒    ▓█   ▀▒████▄     ██▄█▒ ▓██   ▒▓██▒ ██ ▀█   █ ▒██▀ ██▌▓█   ▀ ▓██ ▒ ██▒
▒██░    ▒███  ▒██  ▀█▄  ▓███▄░ ▒████ ░▒██▒▓██  ▀█ ██▒░██   █▌▒███   ▓██ ░▄█ ▒
▒██░    ▒▓█  ▄░██▄▄▄▄██ ▓██ █▄ ░▓█▒  ░░██░▓██▒  ▐▌██▒░▓█▄   ▌▒▓█  ▄ ▒██▀▀█▄  
░██████▒░▒████▒▓█   ▓██▒▒██▒ █▄░▒█░   ░██░▒██░   ▓██░░▒████▓ ░▒████▒░██▓ ▒██▒
░ ▒░▓  ░░░ ▒░ ░▒▒   ▓▒█░▒ ▒▒ ▓▒ ▒ ░   ░▓  ░ ▒░   ▒ ▒  ▒▒▓  ▒ ░░ ▒░ ░░ ▒▓ ░▒▓░
░ ░ ▒  ░ ░ ░  ░ ▒   ▒▒ ░░ ░▒ ▒░ ░      ▒ ░░ ░░   ░ ▒░ ░ ▒  ▒  ░ ░  ░  ░▒ ░ ▒░
  ░ ░      ░    ░   ▒   ░ ░░ ░  ░ ░    ▒ ░   ░   ░ ░  ░ ░  ░    ░     ░░   ░ 
    ░  ░   ░  ░     ░  ░░  ░           ░           ░    ░       ░  ░   ░     
	`

	for _, char := range asciiArt {
		cRed.Print(string(char))
		time.Sleep(4 * time.Millisecond)
	}
	fmt.Println()
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	go AsciiArt(&wg)

	repoPath := scanner.ReceiveRepo()

	wg.Wait()

	var err error
	scanner.CompiledPatterns, err = scanner.GetPatterns()
	if err != nil {
		scanner.ErroColor.Print("[ERRO] ")
		log.Fatalf("Falha ao compilar padrões: %v\n", err)
	}

	scanner.ReadPath(repoPath)
	err = scanner.FinalizeReports(repoPath)
	if err != nil {
		scanner.ErroColor.Print("[ERRO] ")
		log.Fatal(err)
	}
}
