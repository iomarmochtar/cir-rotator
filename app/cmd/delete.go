package cmd

import (
	"fmt"

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
			&cli.StringFlag{
				Name:  "repo-list",
				Usage: "path of file containing repositories that will be deleted, this can be generated from list action",
			},
			&cli.IntFlag{
				Name:  "worker-count",
				Usage: "worker manifest digest parallel deletion count",
				Value: 1,
			},
		}...),
		Action: func(ctx *cli.Context) error {
			cfg, err := initConfig(ctx)
			if err != nil {
				return err
			}

			if pd := cfg.ParallelDeletion(); pd <= 0 {
				return fmt.Errorf("invalid value for worker count: %d, make sure it's more than equal to 1", pd)
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
