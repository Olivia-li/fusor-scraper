package main

func main() {
	db := connectToDB()
	defer db.Close()

	go InsertThreadIntoDB(db)
	
	testSpecificURL("https://fusor.net/board/viewtopic.php?t=15086")

}
