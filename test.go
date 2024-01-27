package main

import (
	"log"
	"strings"
	"sync"

	"github.com/gocolly/colly"
)


// Channel to send data to the database goroutine
var testPostDataChannel = make(chan Thread)

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

	threadId := ""
	if strings.Contains(url, "t=") {
		parts := strings.Split(url, "t=")
		threadId = parts[1]
	}

	// On every thread page call callback
	c.OnHTML("div[class^='post has-profile']", func(e *colly.HTMLElement) {
		// Extract details of the post
		postId := e.Attr("id")
		username := e.ChildText("p.author a.username, p.author a.username-coloured")
		postTime := e.ChildAttr("time", "datetime")
		content := e.ChildText("div[class='content']")

		postData := Thread{
			ThreadId: threadId,
			PostId:   postId,
			Author: username,
			Time: postTime,
			Content:  content,
		}

		// Send the PostData to the database goroutine
		testPostDataChannel <- postData

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
