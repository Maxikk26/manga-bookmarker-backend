package services

import (
	"github.com/gocolly/colly/v2"
	"time"
)

func NewCollector(domainGlob string) (*colly.Collector, error) {
	collector := colly.NewCollector(
		colly.Async(true),
		colly.MaxDepth(1),
		colly.UserAgent("Mozilla/5.0"),
	)

	err := collector.Limit(&colly.LimitRule{
		DomainGlob:  domainGlob,
		Parallelism: 10,
		RandomDelay: 500 * time.Millisecond,
	})
	if err != nil {
		return nil, err
	}

	return collector, nil
}
