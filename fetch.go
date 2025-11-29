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

// should this take a list? then check for each before retrieving from web
func (app *App) FetchDetails() error {
	if app.Config.Day == INT_DEFAULT {
		panic("Scanning for missing days not yet implemented")
	}

	fmt.Printf("Pulling details for AoC %d Day %d ...\n", app.Config.Year, app.Config.Day)

	if err := WriteDescriptionFile(app); err != nil {
		return err
	} else {
		fmt.Printf("Finished fetching %d day-%02d/description.md\n", app.Config.Year, app.Config.Day)
	}

	if err := WriteInputFile(app); err != nil {
		return err
	} else {
		fmt.Printf("Finished fetching %d day-%02d/input.txt\n", app.Config.Year, app.Config.Day)
	}

	return nil
}

func WriteDescriptionFile(app *App) error {
	scraper := CleanScraperWithCookieSet(app)
	url := fmt.Sprintf("https://adventofcode.com/%d/day/%d", app.Config.Year, app.Config.Day)

	var contents []string
	scraper.OnHTML("article", func(e *colly.HTMLElement) {
		content, _ := e.DOM.First().Html()

		markdown, err := htmltomarkdown.ConvertString(content)
		if err != nil {
			log.Fatal(err)
		}

		contents = append(contents, markdown)
	})

	err := scraper.Visit(url)
	if err != nil {
		fmt.Printf("error getting problem spec at %s: %d", url, err)

		return err
	}

	joinedContent := strings.Join(contents, "\n\n\n")

	return WriteStringToFile(getOutputDir(app), "description.md", joinedContent)
}

func WriteInputFile(app *App) error {
	scraper := CleanScraperWithCookieSet(app)
	url := fmt.Sprintf("https://adventofcode.com/%d/day/%d/input", app.Config.Year, app.Config.Day)

	var inputFileContents string
	scraper.OnResponse(func(r *colly.Response) {
		inputFileContents = string(r.Body)
	})

	err := scraper.Visit(url)
	if err != nil {
		fmt.Printf("error getting input at %s: %d", url, err)
		return err
	}

	return WriteStringToFile(getOutputDir(app), "input.txt", inputFileContents)
}

func getOutputDir(app *App) string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return filepath.Join(cwd, fmt.Sprintf("day-%02d", app.Config.Day))
}
