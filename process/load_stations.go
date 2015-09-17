package main

import (
	"encoding/csv"
	"github.com/boltdb/bolt"
	"log"
	"os"
)

func main() {

	log.Println("fred")
	stationlist := os.Getenv("STATIONS")
	log.Println("StationList: ", stationlist)
	log.Println("Loading stations")
	stations := FetchStations(stationlist)
	databasePath := os.Getenv("FETCHDB")

	log.Println("Opening database:", databasePath)
	db, err := bolt.Open(databasePath, 0600, nil)

	if err != nil {
		log.Fatal("Failure Opening database: ", databasePath, err)
	}
	defer db.Close()

	buildabucket(db, "stations")
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

func buildabucket(db *bolt.DB, bucket_name string) {
	db.Update(func(tx *bolt.Tx) error {

		_, err := tx.CreateBucketIfNotExists([]byte(bucket_name)) // working version for now

		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

}
