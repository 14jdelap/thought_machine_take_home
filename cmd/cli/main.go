package main

import (
	"flag"
	"log"
	"os"
	"runtime/debug"
	"strings"

	a "github.com/14jdelap/thought_machine_take_home/internal/auction"
)

func main() {
	filePath := flag.String("path", "inputs/input.txt", "relative path from the the directory where go is being run")
	flag.Parse()

	content, err := os.ReadFile(*filePath)
	if err != nil {
		debug.PrintStack()
		log.Fatal(err)
	}

	splitInputs := strings.Split(string(content), "\n")

	rows := a.ProcessInputs(splitInputs)
	auctionResults := a.HoldAuction(rows)
	a.AnnounceResults(auctionResults)
}
