# GoMock Generator

This tool aims at simplifying and accelerating the generation of mocks in Scalingo projects using
[GoMock](https://github.com/golang/mock/).

This tool can either be used as a CLI or as a Go library.

## CLI

```text
$ gomock_generator -h
NAME:
   GoMock generator - Highly parallelized generator of gomock mocks

USAGE:
   gomock_generator [global options]

VERSION:
   0.1.0

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --base-package value           Project base package. E.g. github.com/Scalingo/go-utils. [$BASE_PACKAGE]
   --mocks-filename value         Filename of the JSON file containing the MockConfiguration. Location of this file is the base package. (default: "mocks.json") [$MOCKS_FILENAME]
   --signatures-filename value    Filename of the signatures cache. Location of this file is the base package. (default: "mocks_sig.json") [$SIGNATURES_FILENAME]
   --concurrent-goroutines value  Concurrent amount of goroutines to generate mock. (default: 4) [$CONCURRENT_GOROUTINES]
   --debug                        Activate debug logs
   --help, -h                     show help
   --version, -v                  print the version
```

## Go Library

The `gomockgenerator` package provides a `GenerateMocks` function. It works along with the
`GenerationConfiguration` and `MocksConfiguration` structures. Comments in the code explain the
purpose of every attribute.

## Installation troubleshooting

## Fix code dependencies

Fix urfave/cli to 934abfb to avoid breaking changes during go build
```
cd $GOPATH/src/github.com/urfave/cli
git checkout master
git pull origin master
git checkout 934abfb
```

## Install external dependencies

```
sudo apt-get install mockgen
sudo apt-get install golang-go.tools
```
