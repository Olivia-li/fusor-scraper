package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/gocolly/colly"
)

func main() {
	// Your main scraper setup and logic here (omitted for brevity)

	// Test the specific URL
	testSpecificURL("https://fusor.net/board/viewtopic.php?t=15086")
}

// Test function for a specific URL
func testSpecificURL(url string) {
	c := colly.NewCollector(
		colly.AllowedDomains("fusor.net", "www.fusor.net"),
	)

	// To avoid visiting the same URL twice
	var visited = struct {
		urls map[string]bool
		sync.RWMutex
	}{urls: make(map[string]bool)}

	// On every thread page call callback
	c.OnHTML("div[class^='post has-profile']", func(e *colly.HTMLElement) {
		// Extract details of the post
		username := e.ChildText("p.author a.username")
		postTime := e.ChildAttr("time", "datetime") 
		content := e.ChildText("div[class='content']")
		fmt.Printf("Content: %q\nUsername: %s\nTime: %s\n\n", content, username, postTime)
	})

	// Handle pagination
	c.OnHTML(`div[class="pagination"] a[class="button"]`, func(e *colly.HTMLElement) {
		pageLink := e.Request.AbsoluteURL(e.Attr("href"))
		
		visited.RLock()
		_, found := visited.urls[pageLink]
		visited.RUnlock()

		if !found {
			visited.Lock()
			visited.urls[pageLink] = true
			visited.Unlock()
			c.Visit(pageLink) // Visit the next page
		}
	})

	// Set error handler
	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	// Start scraping the specific URL
	c.Visit(url)
}
