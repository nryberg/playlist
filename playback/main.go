// main
package main

import (
	//	"fmt"
	_ "./models"
	"./routers"
	"log"
	"net/http"
)

func main() {
	//	fmt.Println("Hello World!")
	routers.Init()
	host := "127.0.0.1:8080"
	log.Println("Opening", host)
	if err := http.ListenAndServe(host, nil); err != nil {
		log.Fatal("Server error:", err)
	}
}
