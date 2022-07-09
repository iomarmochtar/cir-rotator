package cmd

import (
	"github.com/iomarmochtar/cir-rotator/app"
	"github.com/urfave/cli/v2"
)

func DeleteAction() *cli.Command {
	return &cli.Command{
		Name: "delete",
		Flags: append(commonFlags, []cli.Flag{
			&cli.BoolFlag{
				Name:  "dry-run",
				Usage: "just log the action, will not deleting",
				Value: false,
			},
			&cli.StringFlag{
				Name:  "skip-list",
				Usage: "path of file that contains skipping list, will be ignored if matched",
			},
		}...),
		Action: func(ctx *cli.Context) error {
			cfg, err := initConfig(ctx)
			if err != nil {
				return err
			}

			app := app.New(cfg)
			repositories, err := doList(app, ctx)
			if err != nil {
				return err
			}

			if err = app.DeleteRepositories(repositories); err != nil {
				return err
			}

			return nil
		},
	}
}
