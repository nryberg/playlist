// default
package controllers

import (
	//	"fmt"
	"../models"
	"html/template"
	"log"
	"net/http"
	//	"time"
	"strconv"
	"strings"
)

func TracksController(rw http.ResponseWriter, rq *http.Request) {
	log.Println("Controller asking for tracks")
	data, err := models.FetchTracks(5)
	t, err := template.ParseFiles("./views/index.tpl", "./views/track_list.tpl")
	if err != nil {
		log.Println(err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("Executing Tracks")
	t.Execute(rw, data)
}

func ArtistsController(rw http.ResponseWriter, rq *http.Request) {
	data, err := models.FetchArtists(5)
	t, err := template.ParseFiles("./views/index.tpl", "./views/artist_list.tpl")
	if err != nil {
		log.Println(err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(rw, data)
}

func OneArtistController(rw http.ResponseWriter, rq *http.Request) {
	ArtistID := strings.Split(rq.URL.Path, "/")[2]
	ArtistIDint, err := strconv.ParseInt(ArtistID, 10, 64)
	data, err := models.FetchOneArtist(ArtistIDint)
	t, err := template.ParseFiles("./views/index.tpl", "./views/artist.tpl")
	if err != nil {
		log.Println(err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(rw, data)
}

func OneStationController(rw http.ResponseWriter, rq *http.Request) {
	StationID := strings.Split(rq.URL.Path, "/")[2]
	StationIDint, err := strconv.ParseInt(StationID, 10, 64)
	data, err := models.FetchOneStation(StationIDint)
	t, err := template.ParseFiles("./views/index.tpl", "./views/station.tpl")
	if err != nil {
		log.Println(err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(rw, data)
}

/*
func CreateController(rw http.ResponseWriter, rq *http.Request) {
	paste := models.Paste{Title: rq.FormValue("title"), Content: rq.FormValue("content")}
	_, err, now := paste.Add()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	paste.CreatedOn = now
	t, err := template.ParseFiles("src/views/create.tpl")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(rw, paste)
}
*/

/*
func ShowController(rw http.ResponseWriter, rq *http.Request) {
	path := rq.URL.Path
	parms := strings.Split(path, "/")
	rawId := parms[2]
	log.Println("rawId =", rawId)
	/*	id1, _ := strconv.ParseInt(rawId, 10, 64)
		id := int(id1)
		paste := models.GetPaste(id)
	id := models.ToObjectId(rawId)
	paste := models.GetPaste(id)
	t, err := template.ParseFiles("src/views/create.tpl")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(rw, paste)
	//	t.Execute(rw, nil)

}
*/
/*
func main() {
	fmt.Println("Hello World!")
}
*/
