package app

import (
	"github.com/iomarmochtar/cir-rotator/pkg/config"
	"github.com/iomarmochtar/cir-rotator/pkg/registry"
	u "github.com/iomarmochtar/cir-rotator/pkg/usecases"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func doList(c config.IConfig) ([]registry.Repository, error) {
	repositories, err := u.ListRepositories(c.ImageRegistry(), c.IncludeEngine(), c.ExcludeEngine())
	if err != nil {
		return nil, err
	}

	if c.IsOutputTable() {
		printTable(repositories)
	}

	if outputJSON := c.OutputJSONPath(); outputJSON != "" {
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

			if _, err = doList(cfg); err != nil {
				return err
			}

			return nil
		},
	}
}
