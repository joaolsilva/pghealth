package main

import (
	"github.com/joaolsilva/pghealth/cmd"
	"os"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
