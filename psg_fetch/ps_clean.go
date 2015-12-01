package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"time"
)

func main() {
	db, err := SetupDB()
	if err != nil {
		log.Fatal(err)
	}

	queryStmt, err := db.Prepare("SELECT time, stationid, artistid, songid FROM play WHERE stationid = $1 AND songid <> 0 ORDER BY time")

	// queryStmt, err := db.Prepare("SELECT COUNT(*) FROM play WHERE stationid = $1")
	if err != nil {
		log.Fatal(err)
	}
	rows, err := queryStmt.Query("1469")
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var time time.Time
		var stationid int
		var artistid int
		var songid int
		if err := rows.Scan(&time, &stationid, &artistid, &songid); err != nil {
			log.Fatal("Error unloading data: ", err)
		}
		log.Println("Data: ", time, stationid, artistid, songid)
	}

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
