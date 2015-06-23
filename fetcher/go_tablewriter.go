package main

import (
	"encoding/csv"
	"fmt"
	"github.com/olekukonko/tablewriter"
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

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(rawCSVdata[0])

	for index, val := range rawCSVdata {
		if index == 0 {
			table.SetHeader(val)
		} else {
			table.Append(val)
		}
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	table.Render()

}
