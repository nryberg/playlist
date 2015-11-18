package main

import (
	"bufio"
	"database/sql"
	_ "encoding/csv"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

type Station struct {
	Row      int64
	Freq     string
	Location string
	ID       string
}

func main() {
	db, err := SetupDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	filename := "/Users/Nick/Develop/playlist/process/stations.csv"
	stations, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer stations.Close()
	var lines []string
	scanner := bufio.NewScanner(stations)
	// Skip the first line
	_ = scanner.Scan()

	// Run through the rest
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	log.Println("Line Count: ", len(lines))

}

func SetupDB() (*sql.DB, error) {
	username := os.Getenv("DBUSER_WRITE") // "nick" // for dev
	log.Println("Username: ", username)
	pass := os.Getenv("DBUSER_WRITE_PW") // "nick" // for dev
	database := os.Getenv("PLAYLISTDB")
	app_status := os.Getenv("APP_STATUS")
	var connection string
	if app_status == "DEV" {
		connection = fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", username, pass, database)
	} else {
		connection = fmt.Sprintf("user=%s password=%s dbname=%s", username, pass, database)
	}

	log.Println("Opening database:", database)
	db, err := sql.Open("postgres", connection)
	if err != nil {
		log.Fatal(err)
	}

	return db, err

}
