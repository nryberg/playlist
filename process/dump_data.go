package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"os"
	//"io"
	//"encoding/csv"
	//"io/ioutil"
	//"strconv"
	"time"
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
	build_test_bed(10)
	dump_test_bed()
}

func dump_test_bed() {
	db, err := openDB_ReadOnly()
	if err != nil {
		log.Fatal("Failure Opening database: ", err)
	}
	defer db.Close()
	var data Data
	fmt.Println("Line, Time, Artist, ArtistID, SongID, Title, Station, TimeUnix")
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("test"))
		c := b.Cursor()
		i := 0
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if k != nil {
				_ = json.Unmarshal(v, &data)
				timed, _ := time.Parse(time.RFC3339, data.Timestamp)
				time_number := timed.Unix()
				for _, track := range data.Tracks {
					fmt.Printf("%d,%s,%s,%d,%d,%q,%s,%d\n", i, k, track.Artist, track.ArtistID, track.SongID, track.Title, data.StationID, time_number)
				}
				i += 1
			}
		}
		return nil
	})

}

func build_test_bed(limited int) {
	var entries map[string][]byte
	entries = make(map[string][]byte)
	db, err := openDB_ReadOnly()
	if err != nil {
		log.Fatal("Failure Opening database: ", err)
	}
	var data Data
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("tracks"))
		c := b.Cursor()
		k, v := c.First()
		for i := 0; i < limited; i++ {
			if k != nil {
				_ = json.Unmarshal(v, &data)
				if err != nil {
					log.Fatal("Failure : ", err)
				}
				ks := string(k[:])
				entries[ks] = v
				k, v = c.Next()
			}
		}
		return nil
	})
	db.Close()

	db, err = openDB_ReadWrite()
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		_ = tx.DeleteBucket([]byte("test"))
		b, err := tx.CreateBucketIfNotExists([]byte("test"))
		if err != nil {
			log.Fatal("Failure : ", err)
		}
		for k := range entries {
			err = b.Put([]byte(k), entries[k])
			if err != nil {
				log.Println("nope no can do: ", err)
			}
		}
		return nil
	})

}

func int64_to_byte(number int64) []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	_ = binary.PutVarint(buf, number)
	return buf

}

func openDB_ReadOnly() (*bolt.DB, error) {
	databasePath := os.Getenv("TRACKSDB")
	db, err := bolt.Open(databasePath, 0600, &bolt.Options{ReadOnly: true})
	if err != nil {
		log.Fatal(err)
	}

	return db, err
}

func openDB_ReadWrite() (*bolt.DB, error) {
	databasePath := os.Getenv("TRACKSDB")
	db, err := bolt.Open(databasePath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	return db, err
}
