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
	log.Println("Fetching station 1")
	data, stationID, err := FetchAStationNow(stations)
	if err != nil {
		panic(err.Error())
	}

	lastID, err := PushData(data, stationID)

	log.Println("Last rawID :", lastID)

	log.Println("Waiting for next station in 30 secs")
	time.Sleep(30 * time.Second)

	log.Println("Fetching station 2")

	data, stationID, err = FetchAStationNow(stations)
	if err != nil {
		panic(err.Error())
	}
	lastID, err = PushData(data, stationID)

	log.Println("Last rawID :", lastID)
}

func PushData(data string, stationID int64) (int, error) {

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

	stamp := time.Now().Format(time.RFC3339)

	var lastInsertId int
	err = db.QueryRow("INSERT INTO raw(rawtime,stationID, rawdata) VALUES($1,$2,$3) returning rawid;",
		stamp, stationID, data).Scan(&lastInsertId)
	if err != nil {
		panic(err.Error())
	}

	return lastInsertId, err
}

func FetchAStationNow(stations []Station) (string, int64, error) {
	now := time.Now()
	station_number := TimeTwice(now)
	station := stations[station_number]
	log.Println("Fetching station: ", station.Location)
	station_id := station.ID
	data, err := FetchStationData(station_id)
	stationID_int, err := strconv.ParseInt(station_id, 10, 64)
	return data, stationID_int, err

}

func FetchStationData(station_id string) (string, error) {
	url := "http://www.kiisfm.com/services/now_playing.html?streamId=" + station_id + "&limit=12"
	/*
		http://www.kiisfm.com/services/now_playing.html?streamId=
		http://www.1065kissfm.com/services/getNowPlaying.php?streamid=6143&limit=10
		http://www.kissfm.com/services/getNowPlaying.php?streamid=6143&limit=10
	*/
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
