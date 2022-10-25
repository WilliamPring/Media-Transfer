package main

import (
	"fmt"
	"os"

	"github.com/williampring/media-transfer/cmd/terminal"
)

func main() {
	args := os.Args[1:]
	if len(args) < 2 {
		fmt.Println("To few args! To use this please provide two paths")
		fmt.Println("[host path] [remote path]")
		os.Exit(1)
	}
	terminal.Start(args)
}
