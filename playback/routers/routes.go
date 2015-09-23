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
	//http.HandleFunc("/", controllers.IndexController)
	http.HandleFunc("/tracks", controllers.TracksController)
	// http.HandleFunc("/paste/", controllers.ShowController)
}
