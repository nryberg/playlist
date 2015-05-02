package main

import (
	"fmt"
	"github.com/clbanning/mxj"
	"os"
)

/*
type Query struct {
	Title  string `xml:h4`
	Artist string `xml:h5`
}
*/
func main() {
	xmlFile, err := os.Open("./stations/sample.xml")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer xmlFile.Close()
	fmt.Println(xmlFile)

	m, err := mxj.NewMapXmlReader(xmlFile)
	if err != nil {
		fmt.Println("Error unmarshalling reader:", err)
		return
	}

	v, err := m.ValuesForKey("figure")
	if err != nil {
		fmt.Println("Error values for key:", err)
		return
	}
	fmt.Println(v)
	/*
		for _, title := range q.Title {
			fmt.Printf("\t%s\n", title)
		}
	*/
	for key, value := range m {
		fmt.Println("Key:", key, "Value:", value)
	}
}
