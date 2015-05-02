package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/kennygrant/sanitize"
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
				track.Artist = artist
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
	path := "./short_stations/"
	files, err := ioutil.ReadDir(path)
	check(err)
	var past_track Track

	outfile, err := os.Create("./output.csv")
	check(err)

	defer outfile.Close()
	csv_writer := csv.NewWriter(outfile)

	stations := make(map[string][]Track)

	for _, file := range files {
		file_name := file.Name()

		tracks := process_station(file_name, path)
		if file.Name()[22:26] == "KDWB" {
			err = write_tracks(tracks, csv_writer)
			check(err)
		}
		station_name := tracks[0].Station
		if station, ok := stations[station_name]; ok {
			for _, track := range tracks {
				if track != past_track {
					station = append(station, track)
				}
			}
			past_track = tracks[0]
			fmt.Println(past_track)
		} else {
			stations[station_name] = tracks
		}

	}
	csv_writer.Flush()
}
