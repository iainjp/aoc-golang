package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"path/filepath"
)

type ScrapedDetails struct {
	Description string
	Input       string
}

func (app *App) FetchDetails() error {
	if app.Config.Day == INT_DEFAULT {
		return fetchAllMissingDays(app)
	} else {
		return fetchSingleDay(app)
	}
}

func fetchAllMissingDays(app *App) error {
	throttle := time.Tick(2 * time.Second)

	for day := 1; day <= 3; day++ {
		if !inputFileExists(day) {
			app.Config.Day = day
			if err := fetchSingleDay(app); err != nil {
				return err
			}
		} else {
			log.Printf("Input file for Day %d already exists, skipping fetch.\n", day)
		}
		<-throttle
	}
	return nil
}

func fetchSingleDay(app *App) error {
	fmt.Printf("Fetching details for AoC %d Day %d ...\n", app.Config.Year, app.Config.Day)

	details, err := requestSingleDayDetails(app)
	if err != nil {
		return fmt.Errorf("error retrieving details: %d", err)
	}

	if !inputFileExists(app.Config.Day) {
		if err := WriteStringToFile(getOutputDir(app.Config.Day), "input.txt", details.Input); err != nil {
			return err
		}
	}

	if err := WriteStringToFile(getOutputDir(app.Config.Day), "description.md", details.Description); err != nil {
		return err
	}

	return nil
}

func requestSingleDayDetails(app *App) (*ScrapedDetails, error) {
	// fetch description regardless of existing file - to account for Day 2 update
	description, err := app.Scraper.GetDescription(app.Config.Year, app.Config.Day)
	if err != nil {
		return nil, err
	}

	var inputFileContents string

	// only refetch input if not already present
	if !inputFileExists(app.Config.Day) {
		inputFileContents, err = app.Scraper.GetInput(app.Config.Year, app.Config.Day)
		if err != nil {
			return nil, err
		}
	}

	return &ScrapedDetails{
		Description: description,
		Input:       inputFileContents,
	}, err
}

func getOutputDir(day int) string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return filepath.Join(cwd, fmt.Sprintf("day-%02d", day))
}

func inputFileExists(day int) bool {
	filepath := filepath.Join(getOutputDir(day), "input.txt")
	return FileExists(filepath)
}
