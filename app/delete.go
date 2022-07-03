package app

import (
	h "github.com/iomarmochtar/cir-rotator/pkg/helpers"
	u "github.com/iomarmochtar/cir-rotator/pkg/usecases"
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

			repositories, err := doList(cfg)
			if err != nil {
				return err
			}
			var skipList []string
			if skipListPath := ctx.String("skip-list"); skipListPath != "" {
				if skipList, err = h.ReadLines(skipListPath); err != nil {
					return err
				}
			}

			if err = u.DeleteRepositories(cfg.ImageRegistry(), repositories, skipList, ctx.Bool("dry-run")); err != nil {
				return err
			}

			return nil
		},
	}
}
