// routes
package routers

import (
	//	"fmt"
	"../controllers"
	// "github.com/bmizerany/pat"
	// "log"
	"net/http"
)

/*
func main() {
	fmt.Println("Hello World!")
}
*/

func Init() {
	/*
		m := pat.New()
		m.Get("/artist/:artist", http.HandlerFunc(controllers.OneArtistController))
		http.Handle("/", m)
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	*/
	fs := http.FileServer(http.Dir("static"))
	//http.HandleFunc("/", controllers.IndexController)
	http.HandleFunc("/tracks", controllers.TracksController)

	http.HandleFunc("/artists", controllers.ArtistsController)
	http.HandleFunc("/artist/", controllers.OneArtistController)

	http.HandleFunc("/stations", controllers.StationsController)
	http.HandleFunc("/station/", controllers.OneStationController)

	http.Handle("/static/", http.StripPrefix("/static/", fs))
	// http.HandleFunc("/paste/", controllers.ShowController)
}
