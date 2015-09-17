package main

import (
	"encoding/csv"
	"fmt"
	"github.com/crackcomm/go-clitable"
	"os"
)

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

	for index := range rawCSVdata {
		rawCSVdata[index] = rawCSVdata[index][0:6]
	}

	header := rawCSVdata[0]
	row := make(map[string]interface{}, len(header))
	table := clitable.New(header)

	for i := 1; i < 5; i++ {
		for c := range header {
			row[header[c]] = rawCSVdata[i][c]
		}

		table.AddRow(row)
	}
	table.Print()
}
