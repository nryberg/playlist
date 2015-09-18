package main

import (
	_ "encoding/csv"
	"encoding/json"
	"github.com/boltdb/bolt"
	// "io"
	//"io/ioutil"
	"fmt"
	"log"
	"os"
	//"strconv"
)

type Station struct {
	Row      int64
	Freq     string
	Location string
	ID       string
}

func main() {

	stationlist := os.Getenv("STATIONS")
	log.Println("StationList: ", stationlist)

	databasePath := os.Getenv("FETCHDB")
	log.Println("Opening database:", databasePath)
	db, err := bolt.Open(databasePath, 0600, nil)

	if err != nil {
		log.Fatal("Failure Opening database: ", databasePath, err)
	}
	defer db.Close()

	log.Println("Loading stations")
	stations, err := FetchStations(db, "stations")
	log.Println("Station count: ", len(stations))
	fmt.Printf("Row, StationID, Location\n")
	for i, station := range stations {
		fmt.Printf("%d,%s,%s\n", i, station.ID, station.Freq)
	}

}

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
