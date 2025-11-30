package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	colly "github.com/gocolly/colly/v2"
)

type Scraper struct {
	collector     *colly.Collector
	sessionCookie string
}

func NewScraper(sessionCookie string) *Scraper {
	c := colly.NewCollector(
		colly.UserAgent("github.com/iainjp/aoc-golang"),
		colly.AllowedDomains("adventofcode.com"),
		colly.MaxDepth(2),
	)

	c.Limit(&colly.LimitRule{
		DomainRegexp: `adventofcode.com\/\d\d\d\d\/day\/\d+(\/input)?`,
		Parallelism:  1,
		Delay:        3 * time.Second,
		RandomDelay:  2 * time.Second,
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Cookie", "session="+sessionCookie)
	})

	return &Scraper{
		collector:     c,
		sessionCookie: sessionCookie,
	}
}

func (scraper *Scraper) Reset() {
	scraper.collector = scraper.collector.Clone()

	// need to reset the session cookie
	scraper.collector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Cookie", "session="+scraper.sessionCookie)
	})
}

func (scraper *Scraper) GetDescription(year int, day int) (string, error) {
	// hack to remove previously set handlers
	defer scraper.Reset()

	url := fmt.Sprintf("https://adventofcode.com/%d/day/%d", year, day)

	var detailSections []string
	scraper.collector.OnHTML("body main article", func(e *colly.HTMLElement) {
		content, _ := e.DOM.First().Html()

		markdown, err := htmltomarkdown.ConvertString(content)
		if err != nil {
			log.Fatal(err)
		}

		detailSections = append(detailSections, markdown)
	})

	if err := scraper.collector.Visit(url); err != nil {
		return "", err
	}

	return strings.Join(detailSections, "\n\n\n"), nil
}

func (scraper *Scraper) GetInput(year int, day int) (string, error) {
	// hack to remove previously set handlers
	defer scraper.Reset()

	url := fmt.Sprintf("https://adventofcode.com/%d/day/%d/input", year, day)
	var inputFileContents string
	scraper.collector.OnResponse(func(r *colly.Response) {
		fmt.Printf("Visited: %s\n", r.Request.URL.String())
		if strings.Contains(r.Request.URL.String(), "/input") {
			inputFileContents = string(r.Body)
		}
	})

	if err := scraper.collector.Visit(url); err != nil {
		return "", err
	}

	return inputFileContents, nil
}
