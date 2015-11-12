package main

// Grabs Stations data and pushes to PG
import (
	"database/sql"
	"encoding/csv"
	_ "encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

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

func main() {
	stationlist := os.Getenv("STATIONS")
	log.Println("StationList: ", stationlist)
	log.Println("Loading stations")
	stations := FetchStations(stationlist)
	username := "nick" // for dev
	database := os.Getenv("PLAYLISTDB")
	connection := fmt.Sprintf("user=%s dbname=%s sslmode=disable", username, database)
	//db, err := sql.Open("postgres", "user=pqgotest dbname=pqgotest sslmode=verify-full")
	db, err := sql.Open("postgres", connection)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Opening database:", database)

	log.Println("Fetching station 1")
	data, err := FetchAStationNow(stations)
	if err != nil {
		panic(err.Error())
	}
	stamp := time.Now().Format(time.RFC3339)

	var lastInsertId int
	err = db.QueryRow("INSERT INTO raw(rawtime,rawdata) VALUES($1,$2) returning rawid;",
		stamp, data).Scan(&lastInsertId)
	if err != nil {
		panic(err.Error())
	}

	log.Println("Last rawID :", lastInsertId)

}

func FetchAStationNow(stations []Station) (string, error) {
	now := time.Now()
	station_number := TimeTwice(now)
	station := stations[station_number]
	log.Println("Fetching station: ", station.Location)
	station_id := station.ID
	data, err := FetchStationData(station_id)
	return data, err

}

func FetchStationData(station_id string) (string, error) {
	url := "http://www.kiisfm.com/services/now_playing.html?streamId=" + station_id + "&limit=12"

	res, err := http.Get(url)
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)
	body_string := string(body)
	//body, err := ioutil.ReadFile("sample.json")

	if err != nil {
		panic(err.Error())
	}
	return body_string, err

}

func FetchStations(path string) []Station {
	var stations []Station
	station_file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer station_file.Close()

	rdr := csv.NewReader(station_file)
	//rdr.Comma = ','
	// Drop the header row
	_, err = rdr.Read()
	for {
		var station Station
		record, err := rdr.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		station.Row, err = strconv.ParseInt(record[0], 10, 64)
		station.Freq = record[1]
		station.Location = record[2]
		station.ID = record[3]
		stations = append(stations, station)
		if err != nil {
			log.Fatal(err)
		}
	}
	return stations

}

func TimeTwice(t time.Time) int {
	var out float64
	var final int
	working := t.Minute()
	if working >= 30 {
		working -= 30
	}
	out = float64(working)
	if t.Second() > 30 {
		out += .5
	}
	final = int((out * 2))
	return final
}
