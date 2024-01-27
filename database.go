package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv" 
	_ "github.com/lib/pq"      
)


func connectToDB() *sql.DB {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get environment variables
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected to the database!")

	return db
}

func InsertThreadIntoDB(db *sql.DB) {
	for {
		thread := <-testPostDataChannel
		_, err := db.Exec(`INSERT INTO threads (title, content, author, post_time, post_id, thread_id) VALUES ($1, $2, $3, $4, $5, $6)`, thread.Title, thread.Content, thread.Author, thread.Time, thread.PostId, thread.ThreadId)
		if err != nil {
			log.Println("Error inserting data:", err)
		} else {
			fmt.Printf("Post %s data inserted successfully!\n\n", thread.PostId)
		}
	}
}