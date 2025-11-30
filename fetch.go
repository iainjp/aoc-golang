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
		return fetchAllMissingDays(app)
	} else {
		return fetchSingleDay(app)
	}
}

func fetchAllMissingDays(app *App) error {
	for day := 1; day <= 3; day++ {
		if !inputFileExists(day) {
			app.Config.Day = day
			if err := fetchSingleDay(app); err != nil {
				return err
			}
		} else {
			log.Printf("Input file for Day %d already exists, skipping fetch.\n", day)
		}
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

// too much conditionals here - split out
func requestSingleDayDetails(app *App) (*ScrapedDetails, error) {
	url := fmt.Sprintf("https://adventofcode.com/%d/day/%d", app.Config.Year, app.Config.Day)

	// follow link to input file, if input file doesn't already exist
	if !inputFileExists(app.Config.Day) {
		app.Scraper.OnHTML("a[href]", func(h *colly.HTMLElement) {
			link := h.Attr("href")
			if link != "" && strings.Contains(link, `/input`) {
				h.Request.Visit(link)
			}
		})
	} else {
		log.Printf("Input file for Day %d already exists, skipping input fetch.\n", app.Config.Day)
	}

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

	if !inputFileExists(app.Config.Day) {
		app.Scraper.OnResponse(func(r *colly.Response) {
			fmt.Printf("Visited: %s\n", r.Request.URL.String())
			if strings.Contains(r.Request.URL.String(), "/input") {
				inputFileContents = string(r.Body)
			}
		})
	}

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

func inputFileExists(day int) bool {
	filepath := filepath.Join(getOutputDir(day), "input.txt")
	return FileExists(filepath)
}
