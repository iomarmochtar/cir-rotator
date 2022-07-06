package main

import (
	"fmt"
	"os"

	"github.com/iomarmochtar/cir-rotator/app/cmd"
)

func main() {
	a := cmd.New()
	if err := a.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
