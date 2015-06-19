package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/kennygrant/sanitize"
	"gopkg.in/mgo.v2"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"
)

type Track struct {
	Station   string
	TimeStamp time.Time
	Title     string
	Artist    string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func process_station(file string, path string) []Track {
	station, err := os.Open(path + file)
	timestamp_string := file[0:20]
	timestamp, err := time.Parse(time.RFC3339, timestamp_string)
	station_name := file[22:strings.Index(file, ".")]
	check(err)

	defer station.Close()
	var begin_exp = regexp.MustCompile(`<h4`)
	var end_exp = regexp.MustCompile(`<h5`)
	var onair_exp = regexp.MustCompile(`onair-playlist`)
	var is_playlist = false
	scan := bufio.NewScanner(station)
	title := ""
	artist := ""
	var tracks []Track
	track := new(Track)
	for scan.Scan() {
		text := scan.Text()
		if onair_exp.MatchString(text) {
			is_playlist = true
		}

		if is_playlist {
			if begin_exp.MatchString(text) {
				title = sanitize.HTML(text) //strings.Split(text, ">")[1]
				title = strings.Trim(title, " ")
			}
			if end_exp.MatchString(text) {
				artist = sanitize.HTML(text)
				artist = strings.Trim(artist, " ")
				artist = strings.TrimLeft(artist, "by")
				track.Artist = strings.TrimLeft(artist, " ")
				track.Title = title
				track.TimeStamp = timestamp
				track.Station = station_name
				tracks = append(tracks, *track)
			}
		}
	}

	check(err)
	return tracks
}

func write_tracks(tracks []Track, csv_writer *csv.Writer) error {
	var record []string
	var err error
	for _, track := range tracks {
		record = make([]string, 4)
		record[0] = track.Station
		record[1] = track.TimeStamp.String()
		record[2] = track.Artist
		record[3] = track.Title

		err = csv_writer.Write(record)

	}
	return err

}

func main() {
	// Hook up to Mongo for export

	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("music").C("tracks")

	path := "./short_stations/"
	path = "./stations/"
	files, err := ioutil.ReadDir(path)
	check(err)

	stations := make(map[string][]Track)
	station_last_track := make(map[string]Track)

	for _, file := range files {
		file_name := file.Name()

		tracks := process_station(file_name, path)
		if len(tracks) > 0 {
			station_name := file_name[22:strings.Index(file_name, ".")]
			if station, ok := stations[station_name]; ok {
				last_track := station_last_track[station_name]

				for _, track := range tracks {
					if track.Title != last_track.Title {
						station = append(station, track)
						stations[station_name] = station
					} else {
						break
					}

				}

			} else {
				station := append(station, tracks[0])
				stations[station_name] = station
			}
			station_last_track[station_name] = tracks[0]
		}

	}
	/*
		if file.Name()[22:26] == "KDWB" {
			err = write_tracks(tracks, csv_writer)
			check(err)
		}
	*/

	//	fmt.Println(stations)
	keys := make([]string, 0, len(stations))
	for k := range stations {
		keys = append(keys, k)
	}
	outfile, err := os.Create("./output.csv")
	check(err)

	defer outfile.Close()
	// csv_writer := csv.NewWriter(outfile)

	// Drop the test collection
	c.DropCollection()

	for _, key := range keys {
		fmt.Println(key, len(stations[key]))
		tracks := stations[key]
		// err = write_tracks(tracks, csv_writer)
		// check(err)

		for _, track := range tracks {
			err = c.Insert(track)
			check(err)
		}
	}
	// csv_writer.Flush()
	/*
		tracks := stations["KDWB"]
		for _, track := range tracks {
			fmt.Println(track)
		}
	*/
}
