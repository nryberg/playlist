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
			fmt.Printf("key=%s, value=%s\n", k, v)
			err = json.Unmarshal(v, &data)
			if err != nil {
				//fmt.Printf("%T\n%s\n%#v\n", err, err, err)
				return nil
			}
			fmt.Println(data)
			/*
				for i, track := range tracks {
					fmt.Printf("%s - %s (%d) [%s]\n", track.Artist, track.Title, track.SongID, track.TimeStamp)
				}
			*/
			return nil
		})
		/*
			c := b.Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				fmt.Printf("key=%s, value=%s\n", k, v)
			}
		*/
		//p, err = decode(b.Get(k))
		if err != nil {
			return err
		}
		return nil
	})
}
