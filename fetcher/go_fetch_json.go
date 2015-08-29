package main

import (
	"encoding/json"
	// "bufio"
	"fmt"
	"io/ioutil"
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
			TimeStamp string
		} `json:"track"`
	} `json:"tracks"`
}

func main() {
	/*
		url := "http://www.kiisfm.com/services/now_playing.html?streamId=185&limit=12"

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
		fmt.Printf("%d: %s - %s (%d) [%s]\n", i, track.Track.Artist, track.Track.Title, track.Track.SongID, track.Track.TimeStamp)
	}
}
