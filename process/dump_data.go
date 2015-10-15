package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	_ "fmt"
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

func main() {
	build_test_bed(100)
	dump_test_bed()
}

func dump_test_bed() {
	db, err := openDB_ReadOnly()

	if err != nil {
		log.Fatal("Failure Opening database: ", err)
	}
	defer db.Close()

	var entry Entry
	var entries Entries

	// Unload the entries
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
	/*
		fmt.Println("Line, TimeID, StationID, ArtistID, SongID")
		fmt.Printf("%d,%d,%d,%d,%d,%t\n", entry.EntryID, entry.TimeID,
			entry.StationID, entry.ArtistID, entry.SongID, isNewTracks)
	*/
	//TODO Move the dup'ed track info testing down to the pull func
	//			This is after the fact and should be done earlier.  Will
	// 			remove the need to play games with slices
	//
	// This is where the magic happens
	stationEntry := make(map[int64]Entry) // hangs on to the first entry for a station block
	var lastStationID int64               // hangs on the last know Station ID to trip on a station block change
	var lastEntry Entry
	lastStationID = 0

	sort.Sort(entries)
	var isNewTracks bool
	isNewTracks = true
	for i, _ := range entries {
		entry := entries[i]
		if lastStationID == entry.StationID { // working the same station block
			// log.Println("Same : ", entry.SongID, isNewTracks)
			if lastEntry.SongID == entry.SongID { // then we're repeating
				isNewTracks = false
				// log.Println("Repeater: ", entry.SongID, isNewTracks)
			} // is repeated song
		} // is same station block
		if lastStationID != entry.StationID { // change in station block
			// log.Println("New  : ", i, lastStationID, entry.StationID, isNewTracks)
			lastStationID = entry.StationID
			lastEntry = stationEntry[entry.StationID]
			stationEntry[entry.StationID] = entry
			isNewTracks = true
		} // is new station block
		if isNewTracks == true {

		} // print if isNewTrack!
	} // iterate entries

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
				ks := string(k[:])
				chunks[ks] = v
				k, v = c.Next()
			}
		}
		return nil
	})
	db.Close()

	// Walking the dogs

	stationEntry := make(map[int64]Entry) // hangs on to the first entry for a station block
	var lastStationID int64               // hangs on the last know Station ID to trip on a station block change
	var lastEntry Entry
	lastStationID = 0
	var isNewTracks bool

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

			lastStationID = data.StationID
			lastEntry = stationEntry[data.StationID]
			stationEntry[entry.StationID] = entry

			// TODO: Rotate on the correct entry ^^^

			isNewTracks = true
			for _, track := range data.Tracks {
				if track.SongID != 0 {
					nextkey, _ := b.NextSequence()
					entry.EntryID = nextkey
					entry.TimeID = timed.Unix()
					entry.StationID, err = strconv.ParseInt(data.StationID, 10, 64)
					entry.ArtistID = track.ArtistID
					entry.SongID = track.SongID

					if lastEntry.SongID == entry.SongID { // then we're repeating
						isNewTracks = false
						// log.Println("Repeater: ", entry.SongID, isNewTracks)
					} // is repeated song
					enc, err := json.Marshal(entry)
					err = b.Put(uint64_to_byte(nextkey), enc)
					if err != nil {
						log.Println("nope no can do: ", err)
					} // err trap
				} // SongID is not zero, run it
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
