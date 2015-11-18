package main

import (
	"database/sql"
	//"encoding/csv"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	//"io"
	// "io/ioutil"
	"log"
	// "net/http"
	"os"
	// "strconv"
	// "time"

	// "bufio"
	//"github.com/boltdb/bolt"
)

type Data struct {
	Tracks    `json:"tracks"`
	StationID string
	Timestamp string
}

type Tracks []struct {
	Track `json:"track"`
}

type Track struct {
	Artist   string `json:"artistName"`
	ArtistID int64  `json:"thumbplay_artist_id,string"`
	SongID   int64  `json:"thumbplay_song_id,string"`
	Title    string `json:"trackTitle"`
}

type Station struct {
	Row      int64
	Freq     string
	Location string
	ID       string
}

type Artist struct {
	Name     string
	ArtistID int64
	Plays    int64
}

func main() {

	db, err := SetupDB()
	if err != nil {
		log.Fatal(err)
	}

	queryStmt, err := db.Prepare("SELECT rawdata FROM raw")

	artistStmt, err := db.Prepare("INSERT INTO artist(artistid, artistname) SELECT $1, $2 WHERE NOT EXISTS ( SELECT artistid FROM artist WHERE artistid=$1)")

	if err != nil {
		log.Fatal(err)
	}

	rows, err := queryStmt.Query()
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}

	var data Data
	for rows.Next() {
		var rawtext string
		if err := rows.Scan(&rawtext); err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal([]byte(rawtext), &data)
		if err != nil {
			log.Fatal(err)
		}

		for _, track := range data.Tracks {

			_, err = artistStmt.Exec(track.ArtistID, track.Artist)
		}

	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Closing up")
	db.Close()

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
