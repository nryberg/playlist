package main

import (
	"encoding/json"
	// "bufio"
	"fmt"
	"github.com/boltdb/bolt"
	"io/ioutil"
	"log"
	//"net/http"
	"time"
)

type Data struct {
	Tracks []struct {
		Track struct {
			Artist    string `json:"artistName"`
			ArtistID  int64  `json:"thumbplay_artist_id,string"`
			SongID    int64  `json:"thumbplay_song_id,string"`
			Title     string `json:"trackTitle"`
			StationID string
			TimeStamp string
			UNIXTime  int64
		} `json:"track"`
	} `json:"tracks"`
}

func main() {
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	station_id := "185"
	//url := "http://www.kiisfm.com/services/now_playing.html?streamId=" + station_id + "&limit=12"
	/*
		res, err := http.Get(url)
		if err != nil {
			panic(err.Error())
		}

		body, err := ioutil.ReadAll(res.Body)
	*/
	body, err := ioutil.ReadFile("sample.json")

	if err != nil {
		panic(err.Error())
	}

	var data Data

	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Printf("%T\n%s\n%#v\n", err, err, err)
		switch v := err.(type) {
		case *json.SyntaxError:
			fmt.Println(string(body[v.Offset-40 : v.Offset]))
		}
	}

	for i, track := range data.Tracks {
		track.Track.TimeStamp = time.Now().Format(time.RFC3339)
		track.Track.UNIXTime = time.Now().Unix()
		track.Track.StationID = station_id
		err = writetracks(&data, db)
		fmt.Printf("%d: %s - %s (%d) [%s]\n", i, track.Track.Artist, track.Track.Title, track.Track.SongID, track.Track.TimeStamp)
	}
}

func writetracks(data *Data, db *bolt.DB) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("Tracks"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		err = b.Put([]byte("answer"), []byte("42"))
		return nil
	})
	return err
}
