package main

import (
	"fmt"
	"os"

	"github.com/joshmeranda/marina/cmd/operator/app"
)

func main() {
	app := app.App()
	if err := app.Run(os.Args); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
