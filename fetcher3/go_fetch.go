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

	buildabucket(db)
	data.StationID = station_id
	err = writetracks(&data, station_id, db)
	if err != nil {
		panic(err.Error())
	}
}

func writetracks(data *Data, station_id string, db *bolt.DB) error {
	enc, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("tracks"))
		key := []byte(time.Now().Format(time.RFC3339))

		err = b.Put(key, enc)
		return nil
	})
	return err
}

func buildabucket(db *bolt.DB) {
	db.Update(func(tx *bolt.Tx) error {
		/*
			err := tx.DeleteBucket([]byte("tracks"))   // use this for testing - wipe the old one for now.
			_, err = tx.CreateBucket([]byte("tracks")) // use this for testing - wipe the old one for now.
		*/

		_, err := tx.CreateBucketIfNotExists([]byte("tracks")) // working version for now

		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

}
