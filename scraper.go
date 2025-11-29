package main

import (
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

func CleanScraperWithCookieSet(app *App) *colly.Collector {
	scraper := app.Scraper.Clone()
	scraper.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Cookie", "session="+app.SessionCookie)
	})

	return scraper
}
