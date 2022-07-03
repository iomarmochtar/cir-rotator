package main

import (
	"fmt"
	"os"

	"github.com/iomarmochtar/cir-rotator/app"
)

func main() {
	a := app.NewApp()
	if err := a.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
