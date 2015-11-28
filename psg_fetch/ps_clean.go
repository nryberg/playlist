package main

import (
	"database/sql"
	//"encoding/csv"
	"fmt"
	_ "github.com/lib/pq"
	//"io"
	// "io/ioutil"
	"log"
	// "net/http"
	"os"
	// "strconv"

	// "bufio"
	//"github.com/boltdb/bolt"
)

func main() {
	db, err := SetupDB()
	if err != nil {
		log.Fatal(err)
	}

	// queryStmt, err := db.Prepare("SELECT time, stationid, artistid, songid WHERE stationid = $1 ORDER BY time")

	queryStmt, err := db.Prepare("SELECT COUNT(*) FROM play WHERE stationid = $1")
	if err != nil {
		log.Fatal(err)
	}
	rows, err := queryStmt.Query("1469")
	var count int
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}

	rows.Next()
	if err := rows.Scan(&count); err != nil {
		log.Fatal("Error unloading raw: ", err)
	}

	log.Println("Rows: ", count)

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
