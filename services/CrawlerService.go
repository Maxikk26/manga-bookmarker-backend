package services

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"log"
	"net/url"
	"strings"
	"sync"
	"time"
)

func obtainDomainGlob(urlStr string) (string, error) {
	// Parse the URL
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}

	// Extract the host from the URL
	host := parsedURL.Hostname()

	// Remove port if present
	dns := strings.Split(host, ":")[0]

	// Add '*' at the beginning and '/*' at the end
	modifiedDNS := fmt.Sprintf("*%s/*", dns)

	return modifiedDNS, nil
}

func CrawlerService(url string, wg *sync.WaitGroup) {
	start := time.Now()

	c := colly.NewCollector(
		colly.Async(true),              // Enable asynchronous requests
		colly.MaxDepth(1),              // Limit depth to 1 to avoid unnecessary recursion
		colly.UserAgent("Mozilla/5.0"), // Set a common user agent
	)

	domainGlob, err := obtainDomainGlob(url)
	if err != nil {
		log.Fatal(err)
	}

	// Optimize network settings
	err = c.Limit(&colly.LimitRule{
		DomainGlob:  domainGlob,      // Apply limit rule to the specific domain
		Parallelism: 10,              // Increase parallelism
		RandomDelay: 1 * time.Second, // Add random delay to avoid being blocked
	})
	if err != nil {
		log.Fatal("limit error: ", err)
		return
	}

	// div element with class story-info-right to obtain manga title
	c.OnHTML("div.story-info-right h1", func(e *colly.HTMLElement) {
		fmt.Println("title: ", e.Text)
	})

	// span element of div parent to obtain image src of manga
	c.OnHTML("div.story-info-left span.info-image img", func(e *colly.HTMLElement) {
		fmt.Println("image url:", e.Attr("src"))
	})

	// Process only the first li element within ul.row-content-chapter
	c.OnHTML("ul.row-content-chapter li:first-child", func(e *colly.HTMLElement) {
		parts := strings.Fields(e.Text)
		if len(parts) > 1 {
			result := parts[1]
			fmt.Println("result:", result)
		} else {
			fmt.Println("The input string does not contain enough parts.")
		}
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping on url provided
	c.Visit("https://chapmanganato.to/manga-aa951883")

	// Wait for all async tasks to complete
	c.Wait()

	elapsed := time.Since(start) // Calculate the elapsed time
	fmt.Printf("Execution time: %s\n", elapsed)
	defer wg.Done()
}
