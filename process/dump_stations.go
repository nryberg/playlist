package main

import (
	"bytes"
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
)

type Station struct {
	Row      int64
	Freq     string
	Location string
	ID       int64
}

func main() {

	db, err := openDB_ReadOnly()
	if err != nil {
		log.Fatal("Failure Opening database: ", err)
	}
	defer db.Close()
	var station Station
	fmt.Printf("StationID, StationName, Frequency\n")
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("stations"))
		stats := b.Stats()
		log.Println("Bucket Stats - Key Count: ", stats.KeyN)
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			err := json.Unmarshal(v, &station)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%d,%q,%s\n", byte_to_int64(k), station.Location, station.Freq)
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

func byte_to_int64(data []byte) int64 {
	buf := bytes.NewReader(data)
	value, err := binary.ReadVarint(buf)
	if err != nil {
		log.Println("Error loading key val")
	}
	return value

}
