package main

import (
	"fmt"
	"os"

	"github.com/joshmeranda/marina/pkg/cmd/manager"
)

func main() {
	app := manager.App()

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("Error: %s", err)
	}
}
