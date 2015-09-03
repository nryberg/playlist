package main

import (
	"encoding/csv"
	"encoding/json"
	// "bufio"
	"fmt"
	"github.com/boltdb/bolt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Data struct {
	Tracks    `json:"tracks"`
	StationID string
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

	ticker := time.NewTicker(time.Second * 30)

	stations := FetchStations("./stationlist.csv")

	db, err := bolt.Open("tracks_2015_09_02.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	buildabucket(db)
	//	counter := 0
	//	max := 59
	go func() {
		for t := range ticker.C {
			fmt.Println(t)
			station_number := TimeTwice(t)
			station := stations[station_number]
			fmt.Println(t, " - Fetching station: ", station.Location)
			station_id := station.ID
			data := FetchStationData(station_id)
			err = writetracks(&data, station_id, db)
			if err != nil {
				panic(err.Error())
			}
		}
	}()
	time.Sleep(time.Minute * 30)
	ticker.Stop()
	fmt.Println("Ticker stopped")

}

func FetchStationData(station_id string) Data {
	url := "http://www.kiisfm.com/services/now_playing.html?streamId=" + station_id + "&limit=12"

	res, err := http.Get(url)
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)

	//body, err := ioutil.ReadFile("sample.json")

	if err != nil {
		panic(err.Error())
	}

	var data Data

	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Printf("%T\n%s\n%#v\n", err, err, err)
		switch v := err.(type) {
		case *json.SyntaxError:
			fmt.Println(string(body[v.Offset-40 : v.Offset]))
		}
	}

	data.StationID = station_id
	return data

}

func writetracks(data *Data, station_id string, db *bolt.DB) error {
	enc, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("tracks"))
		key := []byte(time.Now().Format(time.RFC3339))

		err = b.Put(key, enc)
		return nil
	})
	return err
}

func buildabucket(db *bolt.DB) {
	db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte("tracks")) // use this for testing - wipe the old one for now.
		/*
			_, err = tx.CreateBucket([]byte("tracks")) // use this for testing - wipe the old one for now.
		*/

		_, err = tx.CreateBucketIfNotExists([]byte("tracks")) // working version for now

		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

}

func FetchStations(path string) []Station {
	var stations []Station
	station_file, err := os.Open("stationlist.csv")
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
	if t.Minute() > 30 {
		out += .5
	}
	final = int((out * 2))
	return final
}
