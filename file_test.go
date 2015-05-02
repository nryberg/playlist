package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	files, err := ioutil.ReadDir("./stations/")
	check(err)
	file := files[0].Name
	fmt.Println(file)
}
