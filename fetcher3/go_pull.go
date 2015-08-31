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
	Tracks
	StationID string
}

type Tracks []struct {
	Track
}

type Track struct {
	Artist    string `json:"artistName"`
	ArtistID  int64  `json:"thumbplay_artist_id,string"`
	SongID    int64  `json:"thumbplay_song_id,string"`
	Title     string `json:"trackTitle"`
	StationID string
	TimeStamp string
	UNIXTime  int64
}

func main() {
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	var data *Data
	err = db.View(func(tx *bolt.Tx) error {
		var err error
		b := tx.Bucket([]byte("tracks"))
		b.ForEach(func(k, v []byte) error {
			//fmt.Printf("key=%s, value=%s\n", k, v)
			err = json.Unmarshal(v, &data)
			if err != nil {
				//fmt.Printf("%T\n%s\n%#v\n", err, err, err)
				return nil
			}
			//fmt.Println(data)
			csv := PrintCSV(data)
			fmt.Println(csv)
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	})
}

func PrintCSV(data *Data) string {
	station_id := data.StationID
	output := station_id + ","
	for i, track := range data.Tracks {
		fmt.Println(track)
		output += fmt.Sprintf("%d,%s,,%s,%d,,%s,\n", i, track.Track.Artist, track.Track.Title, track.Track.SongID, track.Track.TimeStamp)
	}
	return output
}
