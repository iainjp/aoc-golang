package main

import (
	"fmt"
	"log"
	"os"
	"time"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	colly "github.com/gocolly/colly/v2"
)

// should this take a list? then check for each before retrieving from web
func FetchDetails(config *Config) error {
	if config.Day == INT_DEFAULT {
		panic("Scanning for missing days not yet implemented")
	}

	fmt.Printf("Pulling details for AoC %d Day %d ...\n", config.Year, config.Day)

	url := getUrl(config.Year, config.Day)

	fmt.Printf("URL: %s\n", url)

	pageContent, err := getPage(url)
	if err != nil {
		return err
	}

	fmt.Printf("Page Content:\n%s\n", pageContent)
	return nil
}

func getUrl(year int, day int) string {
	return fmt.Sprintf("https://adventofcode.com/%d/day/%d", year, day)
}

func configureScraper() *colly.Collector {
	c := colly.NewCollector(
		colly.UserAgent("github.com/iainjp/aoc-golang by iainpritchard92@gmail.com"),
		colly.AllowedDomains("adventofcode.com"),
		colly.MaxDepth(1),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*adventofcode.com*",
		Parallelism: 1,
		Delay:       1 * time.Second,
	})

	c.OnRequest(func(r *colly.Request) {
		// TODO validate this exists
		sessionCookie := os.Getenv("AOC_SESSION_COOKIE")
		r.Headers.Set("Cookie", "session="+sessionCookie)
	})

	return c
}

func getPage(url string) (string, error) {
	c := configureScraper()

	var pageContent string

	c.OnHTML("article", func(e *colly.HTMLElement) {
		content, _ := e.DOM.First().Html()

		markdown, err := htmltomarkdown.ConvertString(content)
		if err != nil {
			log.Fatal(err)
		}

		pageContent = markdown
	})

	err := c.Visit(url)
	if err != nil {
		return "", err
	}

	return pageContent, nil
}
