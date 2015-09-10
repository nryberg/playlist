package main

import (
	"encoding/json"
	"errors"
	"github.com/boltdb/bolt"
	"log"
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

func main() {
	//db, err := openDatabase("./data/tracks.db")
	db, err := bolt.Open("./data/tracks.db", 0600, nil)
	//	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Database: ", db.Path())
	data, err := fetchTracks(db)
	log.Println(data.Timestamp, data.StationID)

}

func fetchTracks(db *bolt.DB) (Data, error) {
	var data Data
	log.Println("In Fetch: ", db.Path())
	err := errors.New("What happened?")
	err = db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte("tracks"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			err = json.Unmarshal(v, &data)
			if err != nil {
				log.Fatal(err)
			}
			key := string(k[:])
			data.Timestamp = key
			log.Println(key)

		}

		//k, v := c.Last()
		/*
			err = json.Unmarshal(v, &data)
			if err != nil {
				log.Fatal(err)
			}
			key := string(k[:])
			data.Timestamp = key
			log.Println(key)
		*/

		return nil
	})

	return data, err
}

func openDatabase(path string) (bolt.DB, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	return *db, err
}
