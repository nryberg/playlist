package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"os"
	"sort"
	"strconv"
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

type Entry struct {
	EntryID   uint64 `json:"entryID"`
	TimeID    int64  `json:"timeID"`
	StationID int64  `json:"stationID"`
	ArtistID  int64  `json:"artistID"`
	SongID    int64  `json:"songID"`
}

type Entries []Entry

func main() {
	build_test_bed(100)
	dump_test_bed()
}

func dump_test_bed() {
	db, err := openDB_ReadOnly()

	if err != nil {
		log.Fatal("Failure Opening database: ", err)
	}
	defer db.Close()

	var entry Entry
	var entries Entries
	fmt.Println("Line, TimeID, StationID, ArtistID, SongID")
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("test"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if k != nil {
				_ = json.Unmarshal(v, &entry)
				entries = append(entries, entry)
			}
		}
		return nil
	})
	sort.Sort(entries)

	for _, entry := range entries {
		// entry = entries[k]
		fmt.Printf("%d,%d,%d,%d,%d\n", entry.EntryID, entry.TimeID,
			entry.StationID, entry.ArtistID, entry.SongID)
	}

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
	var entry Entry
	err = db.Update(func(tx *bolt.Tx) error {
		_ = tx.DeleteBucket([]byte("test"))
		b, err := tx.CreateBucketIfNotExists([]byte("test"))
		if err != nil {
			log.Fatal("Failure : ", err)
		}
		for k := range entries {

			_ = json.Unmarshal(entries[k], &data)
			timed, _ := time.Parse("2006-01-02T15:04:05-07:00", data.Timestamp)
			for _, track := range data.Tracks {

				nextkey, _ := b.NextSequence()
				entry.EntryID = nextkey
				entry.TimeID = timed.Unix()
				entry.StationID, err = strconv.ParseInt(data.StationID, 10, 64)
				entry.ArtistID = track.ArtistID
				entry.SongID = track.SongID

				enc, err := json.Marshal(entry)
				err = b.Put(uint64_to_byte(nextkey), enc)
				if err != nil {
					log.Println("nope no can do: ", err)
				}
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

func uint64_to_byte(number uint64) []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	_ = binary.PutUvarint(buf, number)
	return buf

}

func byte_to_uint64(val []byte) uint64 {
	buf := bytes.NewBuffer(val)
	value, _ := binary.ReadUvarint(buf)
	return value
}

func byte_to_int64(val []byte) int64 {
	buf := bytes.NewBuffer(val)
	value, _ := binary.ReadVarint(buf)
	return value
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

func (slice Entries) Len() int {
	return len(slice)
}

func (slice Entries) Less(i, j int) bool {
	return slice[i].EntryID < slice[j].EntryID
}

func (slice Entries) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
