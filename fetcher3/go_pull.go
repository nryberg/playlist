package main

import (
	"encoding/json"
	// "bufio"
	"fmt"
	"github.com/boltdb/bolt"
	//"io/ioutil"
	"log"
	//"net/http"
	//"time"
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

func main() {
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	var data *Data
	fmt.Println("Time,Station,LineID,ArtistID,TrackID,Artist,Title")

	err = db.View(func(tx *bolt.Tx) error {
		var err error
		b := tx.Bucket([]byte("tracks"))
		b.ForEach(func(k, v []byte) error {
			//fmt.Printf("key=%s, value=%s\n", k, v)
			//fmt.Printf("key=%s, value=%s\n", k, v)
			err = json.Unmarshal(v, &data)
			if err != nil {
				return nil
			}

			csv := PrintCSV(k, data)
			fmt.Printf(csv)
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	})
}

func PrintCSV(timestamp []byte, data *Data) string {
	//output := "Time,Station,LineID,ArtistID,TrackID,Artist,Title\n"
	output := ""
	timestring := fmt.Sprintf("%s", timestamp)
	station_id := data.StationID
	for i, track := range data.Tracks {
		output += timestring + "," + station_id + ","
		output += fmt.Sprintf("%d,%d,%d,%s,%s\n", i, track.Track.ArtistID, track.Track.SongID, track.Track.Artist, track.Track.Title)
	}
	return output
}
