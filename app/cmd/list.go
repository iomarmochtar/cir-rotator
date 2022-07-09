package cmd

import (
	"fmt"

	"github.com/iomarmochtar/cir-rotator/app"
	"github.com/urfave/cli/v2"
)

func ListAction() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Flags: commonFlags,
		Action: func(ctx *cli.Context) error {
			cfg, err := initConfig(ctx)
			if err != nil {
				return err
			}

			if _, err = doList(app.New(cfg), ctx); err != nil {
				return err
			}

			return nil
		},
		Before: func(ctx *cli.Context) error {
			// must specified the output
			if !ctx.Bool("output-table") && ctx.String("output-json") == "" {
				return fmt.Errorf("must specified one or more output")
			}
			return nil
		},
	}
}
