package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Scalingo/go-utils/gomock_generator/gomockgenerator"
	"github.com/Scalingo/go-utils/logger"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

type app struct {
	config gomockgenerator.GenerationConfiguration
	cli    *cli.App
}

func main() {
	cli.AppHelpTemplate = fmt.Sprintf(`%s
EXAMPLE:

%s --concurrent-goroutine 8 --mocks-filename mymocks.json --signatures-filename sigs.json

Reads the mymocks.json file from the current directory and generates the mocks, 8 goroutines at a time. The signatures of the mocks are stored in sigs.json, in the folder designated by the base package written in mymocks.json.

`, cli.AppHelpTemplate, os.Args[0])

	app := app{
		config: gomockgenerator.GenerationConfiguration{},
		cli:    cli.NewApp(),
	}
	app.cli.Name = "GoMock generator"
	app.cli.Usage = "Highly parallelized generator of gomock mocks"
	app.cli.Version = "0.1.0"
	app.cli.Flags = []cli.Flag{
		cli.StringFlag{Name: "mocks-filepath", Value: "./mocks.json", Usage: "Path to the JSON file containing the MockConfiguration. Location of this file is the base package.", EnvVar: "MOCKS_FILEPATH"},
		cli.StringFlag{Name: "signatures-filename", Value: "mocks_sig.json", Usage: "Filename of the signatures cache. Location of this file is the base package.", EnvVar: "SIGNATURES_FILENAME"},
		cli.IntFlag{Name: "concurrent-goroutines", Value: 4, Usage: "Concurrent amount of goroutines to generate mock.", EnvVar: "CONCURRENT_GOROUTINES"},
		cli.BoolFlag{Name: "debug", Usage: "Activate debug logs"},
	}
	cli.AppHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}
USAGE:
   {{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}
   {{if len .Authors}}
AUTHOR:
   {{range .Authors}}{{ . }}{{end}}
   {{end}}{{if .Commands}}
COMMANDS:
{{range .Commands}}{{if not .HideHelp}}   {{join .Names ", "}}{{ "\t"}}{{.Usage}}{{ "\n" }}{{end}}{{end}}{{end}}{{if .VisibleFlags}}
GLOBAL OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}{{if .Copyright }}
COPYRIGHT:
   {{.Copyright}}
   {{end}}{{if .Version}}
VERSION:
   {{.Version}}
   {{end}}
`
	app.cli.Before = func(c *cli.Context) error {
		app.config.MocksFilePath = c.GlobalString("mocks-filepath")
		app.config.SignaturesFilename = c.GlobalString("signatures-filename")
		app.config.ConcurrentGoroutines = c.GlobalInt("concurrent-goroutines")
		return nil
	}
	app.cli.Action = func(c *cli.Context) error {
		log := logger.Default()
		if c.GlobalBool("debug") {
			log = logger.Default(logger.WithLogLevel(logrus.DebugLevel))
		}
		ctx := logger.ToCtx(context.Background(), log)

		log.WithFields(logrus.Fields{
			"mocks_file_path":       app.config.MocksFilePath,
			"signatures_filename":   app.config.SignaturesFilename,
			"concurrent_goroutines": app.config.ConcurrentGoroutines,
		}).Info("Configuration for this mocks generation")

		rawFile, err := os.Open(app.config.MocksFilePath)
		if err != nil {
			return errors.Wrap(err, "fail to open the mocks file")
		}
		defer rawFile.Close()

		mocksConfiguration := gomockgenerator.MocksConfiguration{}
		err = json.NewDecoder(rawFile).Decode(&mocksConfiguration)
		if err != nil {
			return errors.Wrap(err, "mocks file does not contain valid JSON")
		}

		return gomockgenerator.GenerateMocks(ctx, app.config, mocksConfiguration)
	}

	err := app.cli.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
