package main

import (
	"sync" 
	"fmt"
)

var wg sync.WaitGroup

func main() {
	db := connectToDB()
	defer db.Close()

	wg.Add(1) 
	go func() {
		defer wg.Done() 
		InsertThreadIntoDB(db)
	}()

	wg.Add(1) 
	go func() {
		defer wg.Done() 
		scrapeFusor("https://www.fusor.net/board/")
	}()

	wg.Wait()
	fmt.Println("Entire forum scraped")
}