package main

import (
	"fmt"
	"os"

	"github.com/joshmeranda/marina/cmd/marina-server/app"
)

func main() {
	app := app.App()
	if err := app.Run(os.Args); err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
