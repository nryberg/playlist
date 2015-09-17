package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	infile := "/users/nick/data/financial_insitition/all_2012/all_2012.csv"
	outfile := "../samples/banks.csv"
	sample_size := 10
	line_index := make(map[int]int)

	f, err := os.Open(infile)
	if err != nil {
		fmt.Printf("error opening file: %v\n", err)
		os.Exit(1)
	}

	f_out, err := os.Create(outfile)
	if err != nil {
		fmt.Printf("error opening out file: %v\n", err)
		os.Exit(1)
	}

	defer f_out.Close()

	count, err := lineCounter(f)
	f.Seek(0, 0)

	for counter := 0; counter < sample_size; counter++ {
		sample_line := int(rand.Int63n(int64(count)))
		line_index[sample_line] = counter
	}

	f.Seek(0, 0)
	s := bufio.NewScanner(f)
	counter := 0
	out_count := 0

	for s.Scan() {
		// Write out the header - TODO: make headers optional
		_, ok := line_index[counter]
		if ok == true || counter == 0 {
			out_count += 1
			output := s.Text()
			_, err := f_out.WriteString(output + "\n")
			if err != nil {
				fmt.Printf("error writing to file: %v\n", err)
				os.Exit(1)
			}
		}
		counter += 1
	}
	fmt.Printf("Wrote %d lines\n", out_count)
}

func lineCounter(f *os.File) (int, error) {
	r := bufio.NewReader(f)
	buf := make([]byte, 8196)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return count, err
		}

		count += bytes.Count(buf[:c], lineSep)

		if err == io.EOF {
			break
		}
	}

	return count, nil
}
