package services

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"manga-bookmarker-backend/dtos"
	"manga-bookmarker-backend/models"
	"net/url"
	"regexp"
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
		DomainGlob:  domainGlob,             // Apply limit rule to the specific domain
		Parallelism: 10,                     // Increase parallelism
		RandomDelay: 500 * time.Millisecond, // Add random delay to avoid being blocked
	})
	if err != nil {
		log.Fatal("limit error: ", err)
		return
	}

	// div element with class story-info-right to obtain manga title
	c.OnHTML("div.story-info-right h1", func(e *colly.HTMLElement) {
		data.Name = e.Text
	})

	// span element of div parent to obtain image src of manga
	c.OnHTML("div.story-info-left span.info-image img", func(e *colly.HTMLElement) {
		data.Cover = e.Attr("src")
	})

	// Process only the first li element within ul.row-content-chapter
	c.OnHTML("ul.row-content-chapter li:first-child", func(e *colly.HTMLElement) {
		// Extract the text from the a tag with class chapter-name
		chapterName := strings.ToLower(e.ChildText("a.chapter-name"))

		// Regular expression to match "chapter" followed by a number
		re := regexp.MustCompile(`chapter\s+(\d+)`)
		matches := re.FindStringSubmatch(chapterName)

		data.TotalChapters = "0"
		if len(matches) > 1 {
			data.TotalChapters = matches[1]
		}

		// Extract the text from the span tag with class chapter-time
		chapterTime := e.ChildText("span.chapter-time")

		date, err := ExtractAndParseDateOrTime(chapterTime)
		if err != nil {
			fmt.Println("Error:", err)
		}
		data.LastUpdate = date

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
		DomainGlob:  domainGlob,             // Apply limit rule to the specific domain
		Parallelism: 10,                     // Increase parallelism
		RandomDelay: 500 * time.Millisecond, // Add random delay to avoid being blocked
	})
	if err != nil {
		log.Fatal("limit error: ", err)
		return
	}

	// Process only the first li element within ul.row-content-chapter
	c.OnHTML("ul.row-content-chapter li:first-child", func(e *colly.HTMLElement) {
		// Extract the text from the a tag with class chapter-name
		chapterName := strings.ToLower(e.ChildText("a.chapter-name"))

		// Regular expression to match "chapter" followed by a number
		re := regexp.MustCompile(`chapter\s+(\d+)`)
		matches := re.FindStringSubmatch(chapterName)

		data.TotalChapters = "0"
		if len(matches) > 1 {
			data.TotalChapters = matches[1]
		}

		// Extract the text from the span tag with class chapter-time
		chapterTime := e.ChildText("span.chapter-time")

		date, err := ExtractAndParseDateOrTime(chapterTime)
		if err != nil {
			fmt.Println("Error:", err)
		}
		data.LastUpdate = date

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

func AsyncUpdatesScrapping(url string, manga models.Manga) {
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
		// Extract the text from the a tag with class chapter-name
		chapterName := strings.ToLower(e.ChildText("a.chapter-name"))

		// Regular expression to match "chapter" followed by a number
		re := regexp.MustCompile(`chapter\s+(\d+)`)
		matches := re.FindStringSubmatch(chapterName)

		data.TotalChapters = "0"
		if len(matches) > 1 {
			data.TotalChapters = matches[1]
		}

		// Extract the text from the span tag with class chapter-time
		chapterTime := e.ChildText("span.chapter-time")

		date, err := ExtractAndParseDateOrTime(chapterTime)
		if err != nil {
			fmt.Println("Error:", err)
		}
		data.LastUpdate = date

		if data.TotalChapters != manga.TotalChapters {
			filter := bson.M{"_id": manga.Id}
			err = UpdateManga(data, filter)
			if err != nil {
				fmt.Println("Error updating manga: ", err)
			}
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

func ParseRelativeTime(input string) (time.Time, error) {
	// Split the input string based on " ago"
	parts := strings.Split(input, " ago")
	if len(parts) != 2 {
		return time.Time{}, fmt.Errorf("invalid relative time format")
	}

	// Split the remaining part based on space to get value and unit
	timeParts := strings.Fields(parts[0])
	if len(timeParts) != 2 {
		return time.Time{}, fmt.Errorf("invalid relative time format")
	}

	value, err := strconv.Atoi(timeParts[0])
	if err != nil {
		return time.Time{}, err
	}

	unit := timeParts[1]
	now := time.Now()

	switch unit {
	case "min":
		return now.Add(-time.Minute * time.Duration(value)), nil
	case "hour":
		return now.Add(-time.Hour * time.Duration(value)), nil
	case "day":
		return now.Add(-24 * time.Hour * time.Duration(value)), nil
	default:
		return time.Time{}, fmt.Errorf("unknown time unit: %s", unit)
	}
}

func ParseDate(input string) (time.Time, error) {
	// Split based on ","
	parts := strings.Split(input, ",")
	if len(parts) != 2 {
		return time.Time{}, fmt.Errorf("invalid date format")
	}

	monthDay := strings.TrimSpace(parts[0])
	year := strings.TrimSpace(parts[1])

	// Handle two-digit year
	if len(year) == 2 {
		currentYear := time.Now().Year()
		twoDigitYear := currentYear % 100
		if twoDigitYear > 50 {
			year = fmt.Sprintf("19%s", year)
		} else {
			year = fmt.Sprintf("20%s", year)
		}
	}

	dateString := fmt.Sprintf("%s, %s", monthDay, year)
	layout := "Jan 02, 2006"
	parsedDate, err := time.Parse(layout, dateString)
	if err != nil {
		return time.Time{}, err
	}
	return parsedDate, nil
}

func ExtractAndParseDateOrTime(input string) (time.Time, error) {
	// Check if the input contains "ago"
	if strings.Contains(input, "ago") {
		return ParseRelativeTime(input)
	}

	// Check if the input contains a ","
	if strings.Contains(input, ",") {
		return ParseDate(input)
	}

	return time.Time{}, fmt.Errorf("could not identify date or relative time")
}
