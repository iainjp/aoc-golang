package main


import (
	"fmt"
	"log"
    "os"
    "context"

    "github.com/urfave/cli/v3"
)


func main() {
	cmd := &cli.Command {

		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "year",
				Aliases: []string{"year"},
				Usage:   "AoC year",
				Required: true,
			},
			&cli.IntFlag{
				Name:    "day",
				Aliases: []string{"day"},
				Usage:   "AoC day",
				Required: true,
			},
		},

		Commands : []*cli.Command{
			{
				Name: "pull",
				Usage: "Pull details for given AoC day",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					year := cmd.Int("year")
					day := cmd.Int("day")
					fmt.Printf("Pulling details for AoC %d Day %d ...\n", year, day)
					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
    }
}
