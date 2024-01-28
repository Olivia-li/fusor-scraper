package main

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/gocolly/colly"
)

// Channel to send data to the database goroutine
var postDataChannel = make(chan Thread)

func scrapeFusor(rootURL string) {
	c := colly.NewCollector(
		colly.AllowedDomains("fusor.net", "www.fusor.net"),
	)

	// Visited URLs set
	var visited = struct {
		urls map[string]bool
		sync.RWMutex
	}{urls: make(map[string]bool)}

	// Queue for BFS
	var urlQueue = make(chan string, 1000)
	urlQueue <- rootURL // Enqueue root URL

	// Process thread pages
	c.OnHTML("div[class^='post has-profile']", func(e *colly.HTMLElement) {
		// Extract details of the post
		threadId := ""
		if strings.Contains(e.Request.URL.String(), "t=") {
			parts := strings.Split(e.Request.URL.String(), "t=")
			threadId = parts[1]
		}
		postId := e.Attr("id")
		username := e.ChildText("p.author a.username, p.author a.username-coloured")
		postTime := e.ChildAttr("time", "datetime")
		content := e.ChildText("div[class='content']")

		postData := Thread{
			ThreadId: threadId,
			PostId:   postId,
			Author:   username,
			Time:     postTime,
			Content:  content,
		}

		// Send the PostData to the database goroutine
		postDataChannel <- postData
	})

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		visited.RLock()
		_, found := visited.urls[link]
		visited.RUnlock()

		// Check if the link is useful
		if !found && strings.Contains(link, "fusor.net/board") && !strings.Contains(link, "ucp.php") && !strings.Contains(link, "search.php") && !strings.Contains(link, "donwload") {
			visited.Lock()
			visited.urls[link] = true
			visited.Unlock()
			fmt.Printf("Enqueueing %s\n", link)
			c.Visit(link)
		} else {
			fmt.Printf("Skipping %s\n", link)
		}
	})

	// Set error handler
	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	// Start scraping with BFS
	go func() {
		for {
			select {
			case url := <-urlQueue:
				if url != "" {
					c.Visit(url)
				}
			}
		}
	}()

	// Visit the rootUrl
	c.Visit(rootURL)
}
