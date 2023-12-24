package main

import (
	"fmt"
	"os"

	"github.com/joshmeranda/marina/pkg/cmd/marina"
)

func main() {
	app := marina.App()
	if err := app.Run(os.Args); err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
