package main

import (
	"fmt"
	"log"
	"os"

	"path/filepath"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	colly "github.com/gocolly/colly/v2"
)

// should this take a list? then check for each before retrieving from web
func (app *App) FetchDetails() error {
	config := &app.Config

	if config.Day == INT_DEFAULT {
		panic("Scanning for missing days not yet implemented")
	}

	fmt.Printf("Pulling details for AoC %d Day %d ...\n", config.Year, config.Day)

	if err := WriteDescriptionFile(app); err != nil {
		return err
	}

	if err := WriteInputFile(app); err != nil {
		return err
	}

	return nil
}

func WriteDescriptionFile(app *App) error {
	config := &app.Config
	url := fmt.Sprintf("https://adventofcode.com/%d/day/%d", config.Year, config.Day)

	pageContent, err := getProblemSpec(url, app.Scraper)
	if err != nil {
		fmt.Printf("error retrieving problem spec: %d", err)
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("error retrieving input: %d", err)
		return err
	}

	dir := filepath.Join(cwd, fmt.Sprintf("day-%02d", config.Day))
	return WriteStringToFile(dir, "description.md", pageContent)
}

func WriteInputFile(app *App) error {
	config := &app.Config
	url := fmt.Sprintf("https://adventofcode.com/%d/day/%d/input", config.Year, config.Day)

	inputFileContent, err := getInputFile(url, app.Scraper)
	if err != nil {
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	dir := filepath.Join(cwd, fmt.Sprintf("day-%02d", config.Day))
	return WriteStringToFile(dir, "input.txt", inputFileContent)
}

func getProblemSpec(url string, scraper *colly.Collector) (string, error) {
	var pageContent string

	localScraper := SetSessionCookie(scraper.Clone())

	localScraper.OnHTML("article", func(e *colly.HTMLElement) {
		content, _ := e.DOM.First().Html()

		markdown, err := htmltomarkdown.ConvertString(content)
		if err != nil {
			log.Fatal(err)
		}

		pageContent = markdown
	})

	err := localScraper.Visit(url)
	if err != nil {
		return "", err
	}

	return pageContent, nil
}

func getInputFile(url string, scraper *colly.Collector) (string, error) {
	var inputFileContents string

	localScraper := SetSessionCookie(scraper.Clone())

	localScraper.OnHTML("pre", func(e *colly.HTMLElement) {
		println("Inside pre handler")
		inputFileContents = e.Text
	})

	err := localScraper.Visit(url)
	if err != nil {
		return "", err
	}

	return inputFileContents, nil
}
