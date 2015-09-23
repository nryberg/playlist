// paste
package models

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"log"
	"os"
	//"strings"
	// "time"
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

/*
const DATABASE = "mgopastebin"
const PASTES = "pastes"

func (this Paste) Add() (id bson.ObjectId, err error, createdOn time.Time) {
	session, err1 := GetDB()
	if err1 != nil {
		log.Fatal("Error on database start - Add:", err1)
	}
	collection := session.DB(DATABASE).C(PASTES)
	id = bson.NewObjectId()
	createdOn = time.Now()
	this.Id = id
	this.CreatedOn = createdOn
	err = collection.Insert(&this)
	return
}
*/

func FetchTracks(limit int) (Data, error) {
	databasePath := os.Getenv("TRACKSDB")
	log.Println("Database Path: ", databasePath)
	db, err := bolt.Open(databasePath, 0600, nil)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	var data Data
	err = db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte("tracks"))
		c := b.Cursor()
		k, v := c.First()
		for i := 0; i <= limit; i++ {
			//for k, v := c.First(); k != nil; k, v = c.Next() {
			err = json.Unmarshal(v, &data)
			if err != nil {
				log.Fatal(err)
			}
			key := string(k[:])
			data.Timestamp = key

		}

		return nil
	})

	log.Println(data.Timestamp, data.StationID)
	data.Station, err = FetchOneStation(db, "stations", data.StationID)
	if err != nil {
		log.Fatal(err)
	}
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

/*
func GetTrack(id string) Track {
	session, err1 := GetDB()
	if err1 != nil {
		log.Fatal("Error on database start - GetPaste():", err1)
	}
	collection := session.DB(DATABASE).C(PASTES)
	var paste Paste
	err := collection.FindId(id).One(&paste)
	if err != nil {
		log.Fatal("Error on database get - GetPaste():", err)
	}
	return paste
}

*/

func FetchStations(db *bolt.DB, bucket_name string) ([]Station, error) {
	var stations []Station
	var station Station
	err := db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(bucket_name))
		b.ForEach(func(k, v []byte) error {
			//fmt.Printf("A %s is %s.\n", k, v)
			_ = json.Unmarshal(v, &station)
			stations = append(stations, station)
			return nil
		})
		return nil
	})
	return stations, err
}

func FetchOneStation(db *bolt.DB, bucket_name string, stationName string) (Station, error) {
	var station Station
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket_name))
		v := b.Get([]byte(stationName))
		_ = json.Unmarshal(v, &station)
		return nil
	})
	return station, err
}
