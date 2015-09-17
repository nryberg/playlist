package main

import (
	"bufio"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	filename := "/root/Develop/playlist/fetcher/stations.txt"
	output_file_path := "/root/Develop/playlist/fetcher/station_data/"
	stations, err := os.Open(filename)

	time_now := time.Now().Format(time.RFC3339)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer stations.Close()
	var lines []string
	scanner := bufio.NewScanner(stations)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	var rows [][]string
	var fields []string
	for _, val := range lines {
		//		fmt.Println(val)
		fields = strings.Split(val, ",")

		rows = append(rows, fields)
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Station", "City", "State"})

	for _, v := range rows {
		table.Append(v)
	}

	table.Render()

	for _, row := range rows {
		v := row[0]
		url := "http://m." + v + ".com/"

		fmt.Println("fetching : " + url)
		resp, err := http.Get(url)
		if err != nil {
			// handle error
			fmt.Println(err)
			os.Exit(1)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		outfile := output_file_path + time_now + "--" + v + ".html"
		err = ioutil.WriteFile(outfile, body, 0644)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
