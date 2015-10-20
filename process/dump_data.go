package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"os"
	"sort"
	"strconv"
	"time"
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

type Entry struct {
	EntryID   uint64 `json:"entryID"`
	TimeID    int64  `json:"timeID"`
	StationID int64  `json:"stationID"`
	ArtistID  int64  `json:"artistID"`
	SongID    int64  `json:"songID"`
	NewTrack  bool   `json:"newTrack"`
}

type Entries []Entry

type SongIDs map[int64][]bool

func main() {
	build_test_bed(10000)
	dump_test_bed()
}

func dump_test_bed() {

	// TODO Functional - iterate on a block by block basis
	// Retain the last set of SongID's as a mapped slice
	// if the test songid is in the last batch, fuggabout it
	// Forget about trying to retain any sense of place in the slice.
	// Which sort of makes sense given the drive to randomize slices.
	// Overall assumption is that they won't repeat a song within an hour
	// which is pretty generous.
	// If you know from position in the dump that the song was played earlier
	// and the current song was later, then maybe you'd get away with
	// testing for sub-hour results.  Not really interested in hammering
	// the source with four hits per minute.

	db, err := openDB_ReadOnly()

	if err != nil {
		log.Fatal("Failure Opening database: ", err)
	}
	defer db.Close()

	var entry Entry
	var entries Entries

	// Unload the entries
	fmt.Println("Line, TrackCounter, TimeID, StationID, ArtistID, SongID")
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("test"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if k != nil {
				_ = json.Unmarshal(v, &entry)
				entries = append(entries, entry)
			}
		}
		return nil
	})
	sort.Sort(entries)
	for i, entry := range entries {
		fmt.Printf("%d,%d,%d,%d,%d,%d\n", i, entry.EntryID, entry.TimeID,
			entry.StationID, entry.ArtistID, entry.SongID)
	}

} // end dump_test_bed

func build_test_bed(limited int) {
	var chunks map[string][]byte

	chunks = make(map[string][]byte)
	db, err := openDB_ReadOnly()
	if err != nil {
		log.Fatal("Failure Opening database: ", err)
	}
	var data Data

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("tracks"))
		c := b.Cursor()
		k, v := c.First()
		for i := 0; i < limited; i++ {
			if k != nil {
				_ = json.Unmarshal(v, &data)
				if err != nil {
					log.Fatal("Failure : ", err)
				}
				//fmt.Println(data.StationID)
				ks := string(k[:])
				chunks[ks] = v
				k, v = c.Next()

			}
		}
		return nil
	})
	db.Close()

	var stationHistory map[string]map[int64]bool
	stationHistory = make(map[string]map[int64]bool)

	db, err = openDB_ReadWrite()
	defer db.Close()
	var entry Entry
	err = db.Update(func(tx *bolt.Tx) error {
		_ = tx.DeleteBucket([]byte("test"))
		b, err := tx.CreateBucketIfNotExists([]byte("test"))
		if err != nil {
			log.Fatal("Failure : ", err)
		}

		for k := range chunks {

			_ = json.Unmarshal(chunks[k], &data)
			timed, _ := time.Parse("2006-01-02T15:04:05-07:00", data.Timestamp)
			stationID := data.StationID
			currentStationID, _ := strconv.ParseInt(stationID, 10, 64)

			_, ok := stationHistory[stationID] // test for existing station history, drop unecessary return
			if ok {                            // lets check out what's in the books for this station
				for _, track := range data.Tracks {
					songID := track.SongID

					_, ok := stationHistory[stationID][songID] // don't really care about the value
					if ok {                                    // song already exists and we can drop it
						delete(stationHistory[stationID], songID) // if this doesn't work, we can use this as a history for songs
					} else { // We have a winner! The song doesn't exist, and we can rock n' roll
						stationHistory[stationID][songID] = true
						nextkey, _ := b.NextSequence()
						entry.EntryID = nextkey
						entry.TimeID = timed.Unix()
						entry.StationID = currentStationID
						entry.ArtistID = track.ArtistID
						entry.SongID = track.SongID
						enc, err := json.Marshal(entry)
						if songID != 0 {
							err = b.Put(uint64_to_byte(nextkey), enc)
							if err != nil {
								log.Println("nope no can do: ", err)
							} // err trap
						} // don't store zero track id's
					} // end pushing new song out
				} // iterate tracks
			} else { // never really tracked this station before
				// TODO DRY this up!
				stationHistory[stationID] = make(map[int64]bool)
				for _, track := range data.Tracks {
					songID := track.SongID
					stationHistory[stationID][songID] = true
					nextkey, _ := b.NextSequence()
					entry.EntryID = nextkey
					entry.TimeID = timed.Unix()
					entry.StationID = currentStationID
					entry.ArtistID = track.ArtistID
					entry.SongID = track.SongID
					enc, err := json.Marshal(entry)
					if songID != 0 {
						err = b.Put(uint64_to_byte(nextkey), enc)
						if err != nil {
							log.Println("nope no can do: ", err)
						} // err trap
					} // don't store zero tracks
				} // iterate tracks

			} // iterate tracks in data chunk
		} // iterate data chunk entries
		return nil
	})

}

func int64_to_byte(number int64) []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	_ = binary.PutVarint(buf, number)
	return buf

}

func uint64_to_byte(number uint64) []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	_ = binary.PutUvarint(buf, number)
	return buf

}

func byte_to_uint64(val []byte) uint64 {
	buf := bytes.NewBuffer(val)
	value, _ := binary.ReadUvarint(buf)
	return value
}

func byte_to_int64(val []byte) int64 {
	buf := bytes.NewBuffer(val)
	value, _ := binary.ReadVarint(buf)
	return value
}

func openDB_ReadOnly() (*bolt.DB, error) {
	databasePath := os.Getenv("TRACKSDB")
	db, err := bolt.Open(databasePath, 0600, &bolt.Options{ReadOnly: true})
	if err != nil {
		log.Fatal(err)
	}

	return db, err
}

func openDB_ReadWrite() (*bolt.DB, error) {
	databasePath := os.Getenv("TRACKSDB")
	db, err := bolt.Open(databasePath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	return db, err
}

func (slice Entries) Len() int {
	return len(slice)
}

func (slice Entries) Less(i, j int) bool {
	return slice[i].EntryID < slice[j].EntryID
}

func (slice Entries) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
