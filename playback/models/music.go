// paste
package models

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"github.com/boltdb/bolt"
	"log"
	"os"
	"strconv"
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
	ID       int64
}

type Artist struct {
	Name     string
	ArtistID int64
	Plays    int64
}

func FetchArtists(limit int) ([]Artist, error) {
	db, err := openDB()
	defer db.Close()
	var artist Artist
	var artists []Artist
	err = db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte("artists"))
		c := b.Cursor()
		k, v := c.First()
		for i := 0; i <= limit; i++ {
			//for k, v := c.First(); k != nil; k, v = c.Next() {
			if k != nil {
				artist.ArtistID = byte_to_int64(k)
				artist.Name = string(v)
				artists = append(artists, artist)
			}
			if err != nil {
				log.Fatal(err)
			}

			k, v = c.Next()
		}
		return nil
	})
	return artists, err
}

/*
func FetchStations(limit int) ([]Station, error) {
	db, err := openDB()
	defer db.Close()
	var station Station
	var stations []Station
	log.Println("Fetching Stations")
	err = db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte("stations"))
		c := b.Cursor()
		k, v := c.First()
		for i := 0; i <= limit; i++ {
			//for k, v := c.First(); k != nil; k, v = c.Next() {
			if k != nil {
				err = json.Unmarshal(v, &station)
				if err != nil {
					log.Fatal(err)
				}

				station.ID = string(k)
				stations = append(stations, station)
				log.Println(station.Row, station.ID)
			}
			if err != nil {
				log.Fatal(err)
			}

			k, v = c.Next()
		}
		return nil
	})
	return stations, err
}
*/
func Convert(data []byte) (int64, error) {
	v, err := strconv.ParseUint(string(data), 10, 64)
	if err != nil {
		return 0, err
	}
	return int64(v), nil
}
func FetchOneArtist(id int64) (Artist, error) {
	db, err := openDB()
	defer db.Close()
	var artist Artist
	err = db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte("artists"))
		key := int64_to_byte(id)
		v := b.Get(key)
		if v == nil {
			log.Println("No Key: ", id)
		}
		artist.ArtistID = id
		artist.Name = string(v)
		return nil
	})
	return artist, err
}

func int64_to_byte(number int64) []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	_ = binary.PutVarint(buf, number)
	return buf

}

func byte_to_int64(data []byte) int64 {
	buf := bytes.NewReader(data)
	value, err := binary.ReadVarint(buf)
	if err != nil {
		log.Println("Error loading key val")
	}
	return value

}

func FetchTracks(limit int) (Data, error) {
	db, err := openDB()
	defer db.Close()
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
	db.Close()
	StationID, err := strconv.ParseInt(data.StationID, 10, 64)
	data.Station, err = FetchOneStation(StationID)
	if err != nil {
		log.Fatal(err)
	}
	return data, err
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

func FetchStations(limit int) ([]Station, error) {
	db, err := openDB()
	defer db.Close()
	var stations []Station
	var station Station
	err = db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte("stations"))
		c := b.Cursor()
		k, v := c.First()
		for i := 0; i < limit; i++ {
			if k != nil {
				_ = json.Unmarshal(v, &station)
				stations = append(stations, station)
				k, v = c.Next()
			}
		}
		return nil
	})
	return stations, err
}

func FetchOneStation(id int64) (Station, error) {
	db, err := openDB()
	defer db.Close()
	var station Station
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("stations"))
		v := b.Get(int64_to_byte(id))
		_ = json.Unmarshal(v, &station)
		return nil
	})
	return station, err
}

func openDB() (*bolt.DB, error) {
	databasePath := os.Getenv("TRACKSDB")
	db, err := bolt.Open(databasePath, 0600, &bolt.Options{ReadOnly: true})
	if err != nil {
		log.Fatal(err)
	}

	return db, err
}
