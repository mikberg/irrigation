package main

import (
	"os"

	"github.com/mikberg/irrigation/irrigation/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
