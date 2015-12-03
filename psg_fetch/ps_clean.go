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

	queryTimes, err := db.Prepare("SELECT time FROM play WHERE stationid = $1 GROUP BY time ORDER BY time")
	if err != nil {
		log.Fatal(err)
	}

	//rows, err := queryStmt.Query("1469")
	timeRows, err := queryTimes.Query("1469")
	defer timeRows.Close()
	if err != nil {
		log.Fatal(err)
	}
	var times []time.Time
	var time time.Time
	for timeRows.Next() {
		if err := timeRows.Scan(&time); err != nil {
			log.Fatal("Error unloading data: ", err)
		}
		times = append(times, time)
	}
	timeRows.Close()
	queryText :=
		`UPDATE play 
			SET drop = TRUE
		WHERE stationid = $1 
			AND time = $2
		  AND songid IN (
				SELECT songid 
				FROM play
				WHERE stationid = $1
				AND time = $3)`
	log.Println("Query: ", queryText)
	updateEntries, err := db.Prepare(queryText)
	if err != nil {
		log.Fatal(err)
	} // err trap
	for index, time := range times {
		if index > 0 {
			lastTime := times[index-1]
			_, err := updateEntries.Exec("1469", time, lastTime)
			if err != nil {
				log.Fatal(err)
			}
			/*
				rowCnt, err := res.RowsAffected()
				if err != nil {
					log.Fatal(err)
				}

				log.Println("Rows affected : ", rowCnt)
			*/
		} // index > 0
	} // iterate times
	log.Println("Done updating rows")
} // eof

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
