package main

import (
	"time"

	colly "github.com/gocolly/colly/v2"
)

func ConfigureScraper(sessionCookie string) *colly.Collector {
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

	return c
}
