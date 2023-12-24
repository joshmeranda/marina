package main

import (
	"fmt"
	"os"

	"github.com/joshmeranda/marina/pkg/cmd/server"
)

func main() {
	app := server.App()
	if err := app.Run(os.Args); err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
