package cmd

import (
	"github.com/iomarmochtar/cir-rotator/app"
	"github.com/iomarmochtar/cir-rotator/pkg/registry"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func doList(a *app.App, ctx *cli.Context) ([]registry.Repository, error) {
	repositories, err := a.ListRepositories()
	if err != nil {
		return nil, err
	}

	if ctx.Bool("output-table") {
		printTable(repositories)
	}

	if outputJSON := ctx.String("output-json"); outputJSON != "" {
		if err = dumpToJSON(repositories, outputJSON); err != nil {
			return nil, err
		}

		log.Info().Msgf("json output result written to %s", outputJSON)
	}

	return repositories, nil
}

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
	}
}
