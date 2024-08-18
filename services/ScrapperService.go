package services

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"log"
	"manga-bookmarker-backend/dtos"
	"net/url"
	"strconv"
	"strings"
	"time"
)

//Core services

func MangaScrapping(url string, ch chan<- dtos.MangaScrapperData) {
	start := time.Now()

	var data dtos.MangaScrapperData

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
		data.Name = e.Text
		//fmt.Println("title: ", e.Text)
	})

	// span element of div parent to obtain image src of manga
	c.OnHTML("div.story-info-left span.info-image img", func(e *colly.HTMLElement) {
		//fmt.Println("image url:", e.Attr("src"))
		data.Cover = e.Attr("src")
	})

	// Process only the first li element within ul.row-content-chapter
	c.OnHTML("ul.row-content-chapter li:first-child", func(e *colly.HTMLElement) {
		parts := strings.Fields(e.Text)
		//fmt.Println("parts: ", parts)
		if len(parts) > 1 {
			result := parts[1]
			data.TotalChapters = result

			parsedLastUpdate, err := ExtractAndParseDateOrTime(parts)
			if err != nil {
				fmt.Println(err)
				data.LastUpdate = time.Now()
			}

			fmt.Println(parsedLastUpdate)
			data.LastUpdate = parsedLastUpdate

			//fmt.Println("result:", result)
		} else {
			fmt.Println("The input string does not contain enough parts.")
		}
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping on url provided
	c.Visit(url)

	// Wait for all async tasks to complete
	c.Wait()

	ch <- data

	elapsed := time.Since(start) // Calculate the elapsed time
	fmt.Printf("Execution time: %s\n", elapsed)
	//fmt.Println(fmt.Sprintf("%+v", data))
}

func SyncUpdatesScrapping(url string, ch chan<- dtos.MangaScrapperData) {
	start := time.Now()

	var data dtos.MangaScrapperData

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

	// Process only the first li element within ul.row-content-chapter
	c.OnHTML("ul.row-content-chapter li:first-child", func(e *colly.HTMLElement) {
		parts := strings.Fields(e.Text)
		//fmt.Println("parts: ", parts)
		if len(parts) > 1 {
			result := parts[1]
			data.TotalChapters = result

			parsedLastUpdate, err := ExtractAndParseDateOrTime(parts)
			if err != nil {
				fmt.Println(err)
				data.LastUpdate = time.Now()
			}

			fmt.Println(parsedLastUpdate)
			data.LastUpdate = parsedLastUpdate

			//fmt.Println("result:", result)
		} else {
			fmt.Println("The input string does not contain enough parts.")
		}
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping on url provided
	c.Visit(url)

	// Wait for all async tasks to complete
	c.Wait()

	ch <- data

	elapsed := time.Since(start) // Calculate the elapsed time
	fmt.Printf("Execution time: %s\n", elapsed)
}

//Helpers

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

// ParseRelativeTime parses relative time strings like "1 hour ago".
func ParseRelativeTime(parts []string) (time.Time, error) {
	now := time.Now()

	if len(parts) < 2 {
		return time.Time{}, fmt.Errorf("invalid relative time format")
	}

	quantity, err := strconv.Atoi(parts[0])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid quantity: %w", err)
	}

	unit := parts[1]

	var duration time.Duration
	switch unit {
	case "hour", "hours":
		duration = time.Duration(quantity) * time.Hour
	case "day", "days":
		duration = time.Duration(quantity*24) * time.Hour
	default:
		return time.Time{}, fmt.Errorf("unsupported time unit: %s", unit)
	}

	return now.Add(-duration), nil
}

// ParseDate parses a date string in the format "Aug 12, 24" or "Aug 12,24".
func ParseDate(dateString string) (time.Time, error) {
	parsedTime, err := time.Parse("Jan 2, 06", dateString)
	if err != nil {
		// Try parsing without space after comma
		parsedTime, err = time.Parse("Jan 2,06", dateString)
		if err != nil {
			return time.Time{}, fmt.Errorf("could not parse date: %w", err)
		}
	}
	return parsedTime, nil
}

// ExtractAndParseDateOrTime extracts and parses the correct date or time from the array.
func ExtractAndParseDateOrTime(arr []string) (time.Time, error) {
	if len(arr) < 2 {
		return time.Time{}, fmt.Errorf("array must have at least two elements")
	}

	// Check if the last part is "ago", indicating a relative time
	if arr[len(arr)-1] == "ago" {
		return ParseRelativeTime(arr[len(arr)-3:])
	}

	// Handle date formats by checking the last parts
	for i := len(arr) - 1; i >= 2; i-- {
		if strings.Contains(arr[i], ",") {
			// Construct the date string from the last three elements
			dateString := fmt.Sprintf("%s %s", arr[i-1], arr[i])
			return ParseDate(dateString)
		}
	}

	return time.Time{}, fmt.Errorf("could not identify date or relative time")
}
