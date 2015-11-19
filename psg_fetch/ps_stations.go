package main

import (
	"bufio"
	"database/sql"
	// "encoding/csv"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strings"
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
	filename := os.Getenv("STATIONS") // "nick" // for dev
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

	deleteStmt, err := db.Prepare("DELETE FROM station")
	_, err = deleteStmt.Exec()
	if err != nil {
		log.Fatal("Error running SQL: ", err)
	}

	inserStmt, err := db.Prepare("INSERT INTO station(stationid, stationfreq, stationstate, stationcity) SELECT $1, $2, $3, $4")

	for _, station := range lines {
		data := strings.Split(station, ",")
		log.Println(data)

		inserStmt.Exec(data[0], data[1], data[3][0:len(data[3])-1], data[2][1:len(data[2])])

		if err != nil {
			log.Fatal("Error running insert SQL: ", err)
		}
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
