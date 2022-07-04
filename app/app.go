package app

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/urfave/cli/v2"

	c "github.com/iomarmochtar/cir-rotator/app/config"
	"github.com/iomarmochtar/cir-rotator/pkg/helpers"
	reg "github.com/iomarmochtar/cir-rotator/pkg/registry"
	"github.com/jedib0t/go-pretty/table"
)

const VERSION = "0.1.0"

var (
	commonFlags = []cli.Flag{
		&cli.BoolFlag{
			Name:  "output-table",
			Usage: "show output as table to stdout",
			Value: true,
		},
		&cli.StringFlag{
			Name:  "output-json",
			Usage: "dump result as json file",
		},
		&cli.StringFlag{
			Name:    "basic-auth-user",
			Aliases: []string{"u"},
			Usage:   "basic authentication user",
			EnvVars: []string{"BASIC_AUTH_USER"},
		},
		&cli.StringFlag{
			Name:    "basic-auth-pwd",
			Aliases: []string{"p"},
			Usage:   "basic authentication password",
			EnvVars: []string{"BASIC_AUTH_PWD"},
		},
		&cli.StringFlag{
			Name:     "host",
			Aliases:  []string{"ho"},
			Usage:    "registry host",
			Required: true,
			EnvVars:  []string{"REGISTRY_HOST"},
		},
		&cli.StringFlag{
			Name:    "type",
			Aliases: []string{"t"},
			Usage:   "registry type",
			EnvVars: []string{"REGISTRY_TYPE"},
		},
		&cli.StringFlag{
			Name:    "service-account",
			Aliases: []string{"f"},
			Usage:   "service account file path, it cannot be combined if basic auth args are provided",
			EnvVars: []string{"SA_FILE"},
		},
		&cli.StringSliceFlag{
			Name:    "exclude-filter",
			Aliases: []string{"ef"},
			Usage:   "excluding result",
		},
		&cli.StringSliceFlag{
			Name:    "include-filter",
			Aliases: []string{"if"},
			Usage:   "only process the results of filter",
		},
	}
)

// initConfig create configuration instance from given cmd arguments
func initConfig(ctx *cli.Context) (c.IConfig, error) {
	cfg := &c.Config{
		Debug:              ctx.Bool("debug"),
		OutputTable:        ctx.Bool("output-table"),
		OutputJSON:         ctx.String("output-json"),
		RegUsername:        ctx.String("basic-auth-user"),
		RegPassword:        ctx.String("basic-auth-pwd"),
		ServiceAccountPath: ctx.String("service-account"),
		RegistryHost:       ctx.String("host"),
		RegistryType:       ctx.String("type"),
		ExcludeFilters:     ctx.StringSlice("exclude-filter"),
		IncludeFilters:     ctx.StringSlice("include-filter"),
	}
	if err := cfg.Init(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func printTable(repositories []reg.Repository) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "DIGEST", "REPO", "IMAGE_TAG", "SIZE", "DATE_CREATED", "DATE_UPLOADED"})

	var totalSize uint
	i := 0

	for _, repo := range repositories {
		for _, digest := range repo.Digests {
			i++
			// digest is always prefixed with 'sha256:'
			digestSlug := digest.Name[:27] + "…"

			tagsSlug := strings.Join(digest.Tag, ",")

			if len(tagsSlug) > 30 {
				tagsSlug = tagsSlug[:27] + "…"
			}

			totalSize += digest.ImageSizeBytes
			t.AppendRow([]interface{}{i, digestSlug, repo.Name, tagsSlug, helpers.ByteCountIEC(digest.ImageSizeBytes), digest.Created, digest.Uploaded})
		}
	}

	t.AppendFooter(table.Row{"", "", "", "Total", helpers.ByteCountIEC(totalSize)})
	t.Render()
}

func dumpToJSON(repositories []reg.Repository, jsonPath string) error {
	data, err := json.Marshal(repositories)
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(jsonPath, data, 0600); err != nil {
		return err
	}
	return nil
}

func NewApp() cli.App {
	return cli.App{
		Name:                 "cir-rotator",
		Usage:                "an app for managing container image registry contents",
		Version:              VERSION,
		Compiled:             time.Now(),
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			ListAction(),
			DeleteAction(),
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				Usage:   "Enable debug mode",
				EnvVars: []string{"DEBUG_MODE"},
			},
		},
		Before: func(ctx *cli.Context) error {
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
			if ctx.Bool("debug") {
				zerolog.SetGlobalLevel(zerolog.DebugLevel)
				zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
			} else {
				log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
			}

			return nil
		},
	}
}
