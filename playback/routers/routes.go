// routes
package routers

import (
	//	"fmt"
	"../controllers"
	"net/http"
)

/*
func main() {
	fmt.Println("Hello World!")
}
*/

func Init() {
	fs := http.FileServer(http.Dir("static"))
	//http.HandleFunc("/", controllers.IndexController)
	http.HandleFunc("/tracks", controllers.TracksController)
	http.HandleFunc("/artists", controllers.ArtistsController)
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	// http.HandleFunc("/paste/", controllers.ShowController)
}
