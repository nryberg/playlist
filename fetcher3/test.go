package main

import (
	"fmt"
	"time"
)

func main() {
	n := time.Now()
	fmt.Println(n)
	fmt.Println(n.Second())

	ticker := time.NewTicker(time.Millisecond * 500)

	go func() {
		for t := range ticker.C {
			//			value := t.Second() + t.Nanosecond()
			work := TimeByHalves(t)
			str := fmt.Sprintf("%d.%d : %.2f => %d", t.Second(), t.Nanosecond(), work, TimeTwice(t))
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

func TimeByHalves(t time.Time) float64 {
	var out float64
	working := t.Second()
	if working > 30 {
		working -= 30
	}
	out = float64(working)
	if t.Nanosecond() > 500000000 {
		out += .5
	}

	return out
}

func TimeTwice(t time.Time) int {
	var out float64
	var final int
	working := t.Second()
	if working >= 30 {
		working -= 30
	}
	out = float64(working)
	if t.Nanosecond() > 500000000 {
		out += .5
	}
	final = int((out * 2))
	return final
}
