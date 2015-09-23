package main

import (
	//"encoding/csv"
	"encoding/json"
	"github.com/boltdb/bolt"
	//"io"
	//"io/ioutil"
	"fmt"
	"log"
	"os"
	//"strconv"
)

type Data struct {
	Tracks    `json:"tracks"`
	StationID string
	Timestamp string
	Station
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

	databasePath := os.Getenv("TRACKSDB")

	log.Println("Opening database:", databasePath)
	db, err := bolt.Open(databasePath, 0600, nil)

	if err != nil {
		log.Fatal("Failure Opening database: ", databasePath, err)
	}
	defer db.Close()

	buildabucket(db, "artist_id_name")
	buildabucket(db, "artist_name_id")

	data, err := FetchTracks(db, "tracks", 10)

	for _, entry := range data {
		log.Println(data.StationID)
	}
	//	writeStations(stations, "stations", db)
}

func buildabucket(db *bolt.DB, bucket_name string) {
	db.Update(func(tx *bolt.Tx) error {

		_, err := tx.CreateBucketIfNotExists([]byte(bucket_name)) // working version for now

		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

}

func writeData(key string, value string, bucket_name string, db *bolt.DB) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket_name))
		err := b.Put([]byte(key), []byte(value))
		if err != nil {
			log.Fatal(err)
		}

		return nil
	})
	return err
}

func FetchTracks(db *bolt.DB, bucket_name string, limit int) (Data, error) {
	var data Data
	err := db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte("tracks"))
		c := b.Cursor()
		k, v := c.First()
		for i := 0; i <= limit; i++ {
			//for k, v := c.First(); k != nil; k, v = c.Next() {
			err := json.Unmarshal(v, &data)
			if err != nil {
				log.Fatal(err)
			}
			key := string(k[:])
			data.Timestamp = key

		}

		return nil
	})

	log.Println(data.Timestamp, data.StationID)
	if err != nil {
		log.Fatal(err)
	}
	return data, err
}
