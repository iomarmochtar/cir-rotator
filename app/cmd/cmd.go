package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/iomarmochtar/cir-rotator/app/config"

	"github.com/iomarmochtar/cir-rotator/app"
	"github.com/iomarmochtar/cir-rotator/pkg/helpers"
	reg "github.com/iomarmochtar/cir-rotator/pkg/registry"
	"github.com/jedib0t/go-pretty/table"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/urfave/cli/v2"
)

var (
	// CmdName name of command line
	CmdName = "cir-rotator"
	// Version app version, this will be injected/modified during compilation time
	Version = "0.0.0"
	// BuildHash git commit hash during build process
	BuildHash = "0000000000000000000000000000000000000000"
)

var (
	commonFlags = []cli.Flag{
		&cli.BoolFlag{
			Name:  "output-table",
			Usage: "show output as table to stdout",
		},
		&cli.StringFlag{
			Name:  "output-json",
			Usage: "dump result as json file",
		},
		&cli.BoolFlag{
			Name:    "allow-insecure",
			Usage:   "allow insecure ssl verify",
			EnvVars: []string{"ALLOW_INSECURE_SSL"},
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
		&cli.IntFlag{
			Name:  "worker-count",
			Usage: "http client worker count",
			Value: 1,
		},
	}
)

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

	if err = os.WriteFile(jsonPath, data, 0600); err != nil {
		return err
	}
	return nil
}

func doList(a *app.App, ctx *cli.Context) ([]reg.Repository, error) {
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

// initConfig create configuration instance from given cmd arguments
func initConfig(ctx *cli.Context) (config.IConfig, error) {
	cfg := &config.Config{
		RegUsername:        ctx.String("basic-auth-user"),
		RegPassword:        ctx.String("basic-auth-pwd"),
		ServiceAccountPath: ctx.String("service-account"),
		RegistryHost:       ctx.String("host"),
		RegistryType:       ctx.String("type"),
		SkipListPath:       ctx.String("skip-list"),
		RepoListPath:       ctx.String("repo-list"),
		DryRun:             ctx.Bool("dry-run"),
		ExcludeFilters:     ctx.StringSlice("exclude-filter"),
		IncludeFilters:     ctx.StringSlice("include-filter"),
		AllowInsecure:      ctx.Bool("allow-insecure"),
		WorkerCount:        ctx.Int("worker-count"),
		SkipErrDelete:      ctx.Bool("skip-error"),
	}
	if err := cfg.Init(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func New() cli.App {
	cli.VersionPrinter = func(ctx *cli.Context) {
		_, _ = fmt.Fprintf(ctx.App.Writer, `{"version": "%s", "commit": "%s", "compile_time": "%v"}`,
			ctx.App.Version, BuildHash, ctx.App.Compiled)
	}
	return cli.App{
		Name:                 CmdName,
		Usage:                "an app for managing container image registry contents",
		Version:              Version,
		Compiled:             time.Now(),
		EnableBashCompletion: true,
		Authors: []*cli.Author{
			{
				Name:  "Imam Omar Mochtar",
				Email: "iomarmochtar@gmail.com",
			},
		},
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
