package main

import (
	"fmt"
	"io/ioutil"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	files, err := ioutil.ReadDir("./stations/")
	check(err)
	file := files[0].Name()
	fmt.Println(file)
}
