package main

import (
	"context"
	"fmt"
	"log"
	"os"

	yaml "github.com/goccy/go-yaml"
	"github.com/urfave/cli/v3"
)

type Config struct {
	Year int
	Day  int
}

const AOC_CONFIG_FILE = "aoc.yaml"

func main() {
	cmd := &cli.Command{
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "year",
				Aliases: []string{"year"},
				Usage:   fmt.Sprintf("AoC year (optional if %s config file present).", AOC_CONFIG_FILE),
			},
			&cli.IntFlag{
				Name:    "day",
				Aliases: []string{"day"},
				Usage:   "AoC day (optional). If not provided, will scan for missing days.",
			},
		},

		Commands: []*cli.Command{
			{
				Name:   "fetch",
				Usage:  "Fetch details for given AoC day",
				Action: fetchAction,
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

const INT_DEFAULT = 0

// TODO extract validation logic
func fetchAction(ctx context.Context, cmd *cli.Command) error {
	passedYear := cmd.Int("year")
	passedDay := cmd.Int("day")

	fileConfig := getConfigFromFile()
	if fileConfig.Year == INT_DEFAULT && passedYear == INT_DEFAULT {
		return cli.Exit(fmt.Sprintf("Year must be provided via --year flag or %s config file", AOC_CONFIG_FILE), 1)
	}

	if passedDay == INT_DEFAULT {
		log.Default().Println("--day flag not provided; defaulting to scanning for missing days.")
	}

	var configuredYear int
	if passedYear != INT_DEFAULT {
		configuredYear = passedYear
	} else {
		configuredYear = fileConfig.Year
	}

	config := &Config{
		Year: configuredYear,
		Day:  passedDay,
	}

	FetchDetails(config)
	return nil
}

type ConfigFromFile struct {
	Year int `yaml:"year"`
}

func getConfigFromFile() *ConfigFromFile {
	f, err := os.ReadFile(AOC_CONFIG_FILE)

	if err != nil {
		log.Default().Printf("Could not find %s config file in current directory", AOC_CONFIG_FILE)
		return nil
	}

	var config ConfigFromFile

	if err := yaml.Unmarshal(f, &config); err != nil {
		log.Default().Printf("Could not parse %s config file: %v\n", AOC_CONFIG_FILE, err)
		return nil
	}

	log.Default().Printf("Retrieved %s\n", f)
	return &config
}
