package main

import (
	"strings"
	"time"

	colly "github.com/gocolly/colly/v2"
)

// scraper config shared across instances
func ConfigureScraper(sessionCookie string) *colly.Collector {
	c := colly.NewCollector(
		colly.UserAgent("github.com/iainjp/aoc-golang by iainpritchard92@gmail.com"),
		colly.AllowedDomains("adventofcode.com"),
		colly.MaxDepth(2),
	)

	c.Limit(&colly.LimitRule{
		DomainRegexp: `adventofcode.com\/\d\d\d\d\/day\/\d+(\/input)?`,
		Parallelism:  1,
		Delay:        3 * time.Second,
		RandomDelay:  2 * time.Second,
	})

	// follow link to input file
	c.OnHTML("a[href]", func(h *colly.HTMLElement) {
		link := h.Attr("href")
		if link != "" && strings.Contains(link, `/input`) {
			h.Request.Visit(link)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Cookie", "session="+sessionCookie)
	})

	return c
}
