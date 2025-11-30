package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"path/filepath"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	colly "github.com/gocolly/colly/v2"
)

type ScrapedDetails struct {
	Description string
	Input       string
}

func (app *App) FetchDetails() error {
	if app.Config.Day == INT_DEFAULT {
		panic("Scanning for missing days not yet implemented")
	} else {
		return fetchSingleDay(app)
	}
}

func fetchSingleDay(app *App) error {
	fmt.Printf("Fetching details for AoC %d Day %d ...\n", app.Config.Year, app.Config.Day)

	details, err := requestSingleDayDetails(app)
	if err != nil {
		return fmt.Errorf("error retrieving details: %d", err)
	}

	if err := WriteStringToFile(getOutputDir(app.Config.Day), "input.txt", details.Input); err != nil {
		return err
	}

	if err := WriteStringToFile(getOutputDir(app.Config.Day), "description.md", details.Description); err != nil {
		return err
	}

	return nil
}

func requestSingleDayDetails(app *App) (*ScrapedDetails, error) {
	url := fmt.Sprintf("https://adventofcode.com/%d/day/%d", app.Config.Year, app.Config.Day)

	var detailSections []string
	app.Scraper.OnHTML("body main article", func(e *colly.HTMLElement) {
		content, _ := e.DOM.First().Html()

		markdown, err := htmltomarkdown.ConvertString(content)
		if err != nil {
			log.Fatal(err)
		}

		detailSections = append(detailSections, markdown)
	})

	var inputFileContents string
	app.Scraper.OnResponse(func(r *colly.Response) {
		fmt.Printf("Visited: %s\n", r.Request.URL.String())
		if strings.Contains(r.Request.URL.String(), "/input") {
			inputFileContents = string(r.Body)
		}
	})

	err := app.Scraper.Visit(url)
	return &ScrapedDetails{
		Description: strings.Join(detailSections, "\n\n\n"),
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
