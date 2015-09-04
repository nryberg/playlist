package main

import (
	"fmt"
	"time"
)

func main() {
	n := time.Now()
	fmt.Println(n)
	fmt.Println(n.Second())

	ticker := time.NewTicker(time.Second)

	go func() {
		for t := range ticker.C {
			//			value := t.Second() + t.Nanosecond()
			str := fmt.Sprintf("%v :: %d.%d : %d", t, t.Minute(), t.Second(), TimeTwice(t))
			fmt.Println(str)
		}
	}()
	time.Sleep(time.Second * 75)
	ticker.Stop()
	fmt.Println("Ticker stopped")
	max := 30.0
	for i := 0.0; i < max; {
		//fmt.Printf("%3.0f  :  %3.2f\n", i, i*2)
		i += .5
	}
}

func TimeTwice(t time.Time) int {
	var out float64
	var final int
	working := t.Minute()
	if working >= 30 {
		working -= 30
	}
	out = float64(working)
	//if t.Nanosecond() > 500000000 {
	if t.Second() > 30 {
		out += .5
	}
	final = int((out * 2))
	return final
}
