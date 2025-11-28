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
		Commands : []*cli.Command{
			{
				Name: "pull",
				Usage: "Pull details for given AoC day",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					fmt.Println("Pulling ...")
					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
    }
}
