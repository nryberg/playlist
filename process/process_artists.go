package main

import (
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

type Artist struct {
	Name     string
	ArtistID int64
	Plays    int64
}

type Data struct {
	Tracks    `json:"tracks"`
	StationID string
	Timestamp string
	Station
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

type Station struct {
	Row      int64
	Freq     string
	Location string
	ID       int64
}

func main() {

	// Test
	databasePath := os.Getenv("TRACKSDB")

	// Prod
	//databasePath = os.Getenv("FETCHDB")

	log.Println("Opening database:", databasePath)
	db, err := bolt.Open(databasePath, 0600, nil)

	if err != nil {
		log.Fatal("Failure Opening database: ", databasePath, err)
	}
	defer db.Close()

	buildabucket(db, "artists")
	buildabucket(db, "artist_name_id")

	artists, err := FetchAllTrackArtists(db, "tracks")
	log.Println("Entry Count: ", len(data))
	var artist Artist

	data_count := 0
	artist_count := 0
	for _, datum := range data {
		data_count += 1
		for _, track := range datum.Tracks {
			err := db.Update(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte("artists"))
				artist_count += 1
				// TODO: check and see if you've already got this artist -
				artist.ArtistID = track.ArtistID
				artist.Name = track.Artist
				enc, err := json.Marshal(artist)
				if err != nil {
					log.Panic("Artist isn't fitting into struct")
				}
				key := int64_to_byte(track.ArtistID)
				err = b.Put(key, enc)
				if err != nil {
					log.Fatal(err)
				}
				return nil
			})
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	log.Println("Processed data chunks: ", data_count)
	log.Println("Processed artist entries: ", artist_count)

}

func int64_to_byte(number int64) []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	_ = binary.PutVarint(buf, number)
	return buf

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

func writeString(key string, value string, bucket_name string, db *bolt.DB) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket_name))
		err := b.Put([]byte(key), []byte(value))
		if err != nil {
			log.Fatal(err)
		}

		return nil
	})
	return err
}

func FetchTracks(db *bolt.DB, bucket_name string, limit int) ([]Artists, error) {
	var datum Data
	var artists []Artist
	var artist Artist
	err := db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte("tracks"))
		c := b.Cursor()
		//for k, v := c.First(); k != nil; k, v = c.Next() {

		k, v := c.First()
		for i := 0; i <= limit; i++ {
			if k != nil {
				err := json.Unmarshal(v, &datum)
				if err != nil {
					log.Fatal(err)
				}
				for _, a := range datum.Tracks {
					artist.Name = a.Artist
					artist.ArtistID = a.ArtistID
					artists = append(artists, artist)
				} // Pull Artist info
			} // make sure k isn't nil
			k, v = c.Next()
		} // iterate through limited results

		return nil
	})

	log.Println("artist count: ", len(artists))
	if err != nil {
		log.Fatal(err)
	}
	return artists, err

}

func FetchAllTrackArtists(db *bolt.DB, bucket_name string) ([]Artist, error) {
	var datum Data
	var artists []Artist
	var artist Artist
	err := db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte("tracks"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if k != nil {
				err := json.Unmarshal(v, &datum)
				if err != nil {
					log.Fatal(err)
				}

				for _, a := range datum.Tracks {
					artist.Name = a.Artist
					artist.ArtistID = a.ArtistID
					artists = append(artists, artist)

				}
			}
		}

		return nil
	})

	log.Println("artist count: ", len(artists))
	if err != nil {
		log.Fatal(err)
	}
	return artists, err
}
func FetchAllTracks(db *bolt.DB, bucket_name string) ([]Data, error) {
	var data []Data
	var datum Data
	err := db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte("tracks"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if k != nil {
				err := json.Unmarshal(v, &datum)
				if err != nil {
					log.Fatal(err)
				}
				key := string(k[:])
				datum.Timestamp = key
				data = append(data, datum)
			}
		}

		return nil
	})

	log.Println(datum.Timestamp, datum.StationID)
	if err != nil {
		log.Fatal(err)
	}
	return data, err
}
