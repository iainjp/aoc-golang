package main

import (
	"context"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

type Config struct {
	Year int
	Day  int
}

func main() {
	cmd := &cli.Command{
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:     "year",
				Aliases:  []string{"year"},
				Usage:    "AoC year",
				Required: true,
			},
			&cli.IntFlag{
				Name:     "day",
				Aliases:  []string{"day"},
				Usage:    "AoC day",
				Required: true,
			},
		},

		Commands: []*cli.Command{
			{
				Name:  "fetch",
				Usage: "Fetch details for given AoC day",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					config := Config{Year: cmd.Int("year"), Day: cmd.Int("day")}
					FetchDetails(config)
					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
