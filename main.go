package main

import (
	"fmt"
	"os"
)

func main() {
	var command string
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	switch command {
	case "server":
		runServer()
	case "populate-codes":
		populateCodes()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", command)
		os.Exit(1)
	}
}
