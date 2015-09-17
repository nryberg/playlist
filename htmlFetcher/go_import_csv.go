package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

type TestRecord struct {
	Year string
	Name string
}

func main() {
	filename := "../samples/banks.csv"
	csvfile, err := os.Open(filename)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer csvfile.Close()

	reader := csv.NewReader(csvfile)

	reader.FieldsPerRecord = -1

	rawCSVdata, err := reader.ReadAll()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// sanity check, display to standard output
	// for _, each := range rawCSVdata {
	for i := 0; i < 20; i++ {
		each := rawCSVdata[i]
		// fmt.Println(rawCSVdata[i][0])
		fmt.Printf("year: %s and name : %s \n", each[0], each[4])
	}

	// now, safe to move raw CSV data to struct

	var oneRecord TestRecord

	var allRecords []TestRecord

	for _, each := range rawCSVdata {
		oneRecord.Year = each[0]
		oneRecord.Name = each[4]
		allRecords = append(allRecords, oneRecord)
	}

	// second sanity check, dump out allRecords and see if
	// individual record can be accessible
	// fmt.Println(allRecords)
	fmt.Println(allRecords[2].Year)
	fmt.Println(allRecords[2].Name)

}
