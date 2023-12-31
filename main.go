package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "source",
				Aliases:  []string{"s"},
				Usage:    "Source commit hash",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "dest",
				Aliases: []string{"d"},
				Value:   "HEAD",
				Usage:   "Dest commit hash / HEAD",
			},
			&cli.StringFlag{
				Name:    "repo",
				Aliases: []string{"r"},
				Value:   "./",
				Usage:   "Git repository path",
			},
		},
		Action: func(ctx *cli.Context) error {
			var configMap ConfigMap
			configMap.Parse()
			diffFiles := getDiff(ctx.String("repo"), ctx.String("source"), ctx.String("dest"))

			for configKey, config := range configMap {
				fmt.Printf("evaluating ")
				PrintYellow("%s", configKey)
				if !config.IsMatch(diffFiles) {
					PrintGreen(" PASSED!\n")
					continue
				}

				errString := config.GetErrorString()
				if errString != "" {
					PrintRed(" FAIL!\n")
					PrintYellow("%s\n", errString)
					os.Exit(1)
				}

				PrintGreen(" PASSED!\n")
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
