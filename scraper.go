package main

import (
	"os"
	"time"

	colly "github.com/gocolly/colly/v2"
)

// scraper config shared across instances
func ConfigureScraper() *colly.Collector {
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

	return c
}

func SetSessionCookie(scraper *colly.Collector) *colly.Collector {
	scraper.OnRequest(func(r *colly.Request) {
		// TODO validate this exists
		sessionCookie := os.Getenv("AOC_SESSION_COOKIE")
		r.Headers.Set("Cookie", "session="+sessionCookie)
	})

	return scraper
}
