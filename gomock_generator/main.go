package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"

	"github.com/Scalingo/go-utils/errors/v3"
	"github.com/Scalingo/go-utils/gomock_generator/gomockgenerator"
	"github.com/Scalingo/go-utils/logger"
)

var version = "1.4.2"

func main() {
	cli.RootCommandHelpTemplate = fmt.Sprintf(`%s
	EXAMPLE:

	%s --concurrent-goroutine 8 --mocks-filename mymocks.json --signatures-filename sigs.json

	Reads the mymocks.json file from the current directory and generates the mocks, 8 goroutines at a time. The signatures of the mocks are stored in sigs.json, in the folder designated by the base package written in mymocks.json.

	`, cli.RootCommandHelpTemplate, os.Args[0])

	cfgGeneration := gomockgenerator.GenerationConfiguration{}

	cmd := cli.Command{
		Name:      "GoMock generator",
		UsageText: "gomock_generator",
		Usage:     "Highly parallelized generator of gomock mocks",
		Version:   version,
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "mocks-filepath", Value: "./mocks.json", Usage: "Path to the JSON file containing the MockConfiguration. Location of this file is the base package.", Sources: cli.EnvVars("MOCKS_FILEPATH")},
			&cli.StringFlag{Name: "signatures-filename", Value: "mocks_sig.json", Usage: "Filename of the signatures cache. Location of this file is the base package.", Sources: cli.EnvVars("SIGNATURES_FILENAME")},
			&cli.IntFlag{Name: "concurrent-goroutines", Value: 4, Usage: "Concurrent amount of goroutines to generate mock.", Sources: cli.EnvVars("CONCURRENT_GOROUTINES")},
			&cli.BoolFlag{Name: "debug", Usage: "Activate debug logs"},
		},
		CustomRootCommandHelpTemplate: `NAME:
   {{.Name}} - {{.Usage}}
USAGE:
   {{.UsageText}} {{if .VisibleFlags}}[global options]{{end}}
   {{if len .Authors}}
AUTHOR:
   {{range .Authors}}{{ . }}{{end}}
   {{end}}{{if .Commands}}
COMMANDS:
{{range .Commands}}   {{join .Names ", "}}{{ "\t"}}{{.Usage}}{{ "\n" }}{{end}}{{end}}{{if .VisibleFlags}}
GLOBAL OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}{{if .Copyright }}
COPYRIGHT:
   {{.Copyright}}
   {{end}}{{if .Version}}
VERSION:
   {{.Version}}
   {{end}}
`,
		Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
			cfgGeneration.MocksFilePath = cmd.String("mocks-filepath")
			cfgGeneration.SignaturesFilename = cmd.String("signatures-filename")
			cfgGeneration.ConcurrentGoroutines = cmd.Int("concurrent-goroutines")
			return ctx, nil
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			err := validateBinaryDeps(ctx)
			if err != nil {
				return err
			}

			log := logger.Default()
			if cmd.Bool("debug") {
				log = logger.Default(logger.WithLogLevel(logrus.DebugLevel))
			}
			ctx = logger.ToCtx(ctx, log)

			log.WithFields(logrus.Fields{
				"mocks_file_path":       cfgGeneration.MocksFilePath,
				"signatures_filename":   cfgGeneration.SignaturesFilename,
				"concurrent_goroutines": cfgGeneration.ConcurrentGoroutines,
			}).Info("Configuration for this mocks generation")

			rawFile, err := os.Open(cfgGeneration.MocksFilePath)
			if err != nil {
				return errors.Wrap(ctx, err, "fail to open the mocks file")
			}
			defer rawFile.Close()

			mocksConfiguration := gomockgenerator.MocksConfiguration{}
			err = json.NewDecoder(rawFile).Decode(&mocksConfiguration)
			if err != nil {
				return errors.Wrap(ctx, err, "mocks file does not contain valid JSON")
			}

			return gomockgenerator.GenerateMocks(ctx, cfgGeneration, mocksConfiguration)
		},
	}

	err := cmd.Run(context.Background(), os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func validateBinaryDeps(ctx context.Context) error {
	binaries := []struct {
		Executable string
		Package    string
	}{
		{
			Executable: "goimports",
			Package:    "golang.org/x/tools/cmd/goimports",
		},
		{
			Executable: "mockgen",
			Package:    "go.uber.org/mock/mockgen",
		},
	}
	for _, binary := range binaries {
		_, err := exec.LookPath(binary.Executable)
		if err != nil {
			fmt.Fprintf(
				os.Stderr,
				"Executable '%s' not found. Not in $PATH or not installed. Attempt to install.\n\nRunning 'go get %v'...\n\n",
				binary.Executable,
				binary.Package,
			)
			cmd := exec.Command("go", "get", binary.Package)
			err := cmd.Run()
			if err != nil {
				output, outputErr := cmd.CombinedOutput()
				if outputErr != nil {
					return errors.Wrapf(ctx,
						err,
						"run 'go get %v', fail to get command output, error: \n\n%v\n",
						binary.Package, outputErr,
					)
				} else {
					return errors.Wrapf(ctx,
						err,
						"run 'go get %v', output: \n\n%v\n",
						binary.Package, string(output),
					)
				}
			}
			_, err = exec.LookPath(binary.Executable)
			if err != nil {
				return errors.Wrapf(ctx,
					err,
					"find '%s' binary after installation, $GOPATH/bin probably not in path, update your shell (bash/zsh) configuration",
					binary.Executable,
				)
			}
		}
	}
	return nil
}
