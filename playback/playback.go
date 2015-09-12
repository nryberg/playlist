package main

import (
	"encoding/json"
	"errors"
	"github.com/boltdb/bolt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
)

type Data struct {
	Tracks    `json:"tracks"`
	StationID string
	Timestamp string
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
	fs := http.FileServer(http.Dir("static"))
	databasePath := os.Getenv("TRACKSDB")
	log.Println("Database Path: ", databasePath)
	db, err := bolt.Open(databasePath, 0600, nil)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	data, err := fetchTracks(db, 5)
	log.Println(data.Timestamp, data.StationID)

	//http.Handle("/", fs)
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", serveTemplate)

	log.Println("Listening...")
	http.ListenAndServe(":3000", nil)
}

func openDatabase(path string) (bolt.DB, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	return *db, err
}

func fetchTracks(db *bolt.DB, limit int) (Data, error) {
	var data Data
	err := errors.New("What happened?")
	err = db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte("tracks"))
		c := b.Cursor()
		k, v := c.First()
		for i := 0; i <= limit; i++ {
			//for k, v := c.First(); k != nil; k, v = c.Next() {
			err = json.Unmarshal(v, &data)
			if err != nil {
				log.Fatal(err)
			}
			key := string(k[:])
			data.Timestamp = key

		}

		return nil
	})

	return data, err
}
func serveTemplate(w http.ResponseWriter, r *http.Request) {
	lp := path.Join("templates", "layout.html")
	fp := path.Join("templates", r.URL.Path)

	// Return a 404 if the template doesn't exist
	info, err := os.Stat(fp)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
	}

	// Return a 404 if the request is for a directory
	if info.IsDir() {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles(lp, fp)
	if err != nil {
		// Log the detailed error
		log.Println(err.Error())
		// Return a generic "Internal Server Error" message
		http.Error(w, http.StatusText(500), 500)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "layout", nil); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}
