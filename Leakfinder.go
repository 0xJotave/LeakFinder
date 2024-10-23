package main

import (
	"LeakFinder/scanner"
	"fmt"
	"sync"
	"time"

	"github.com/fatih/color"
)

var cRed = color.New(color.FgHiRed).Add(color.Bold)

const asciiArt = `
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

func AsciiArt(wg *sync.WaitGroup) {
	defer wg.Done()

	for _, char := range asciiArt {
		cRed.Print(string(char))
		time.Sleep(4 * time.Millisecond)
	}
	fmt.Println()
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	var artDisplayed bool

	if !artDisplayed {
		go AsciiArt(&wg)
		artDisplayed = true
	}

	repoPath := scanner.ReceiveRepo()

	wg.Wait()

	scanner.CompiledPatterns, _ = scanner.GetPatterns()

	if err := scanner.ReadPath(repoPath); err != nil {
		scanner.HandleError(err.Error(), "")
		return
	}

	if err := scanner.FinalizeReports(repoPath); err != nil {
		scanner.HandleError(err.Error(), "")
	}
}
