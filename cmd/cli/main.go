package main

import (
	"flag"
	"log"
	"os"
	"runtime/debug"

	"github.com/14jdelap/thought_machine_take_home/internal"
)

func main() {
	filePath := flag.String("path", "input.txt", "relative path from the top-level directory to file with auction inputs")
	flag.Parse()

	content, err := os.ReadFile("inputs/" + *filePath)
	if err != nil {
		debug.PrintStack()
		log.Fatal(err)
	}

	house := &internal.AuctionHouse{}
	house.ProcessInputs(content)
	house.HoldAuction()
	house.AnnounceResults()
}
