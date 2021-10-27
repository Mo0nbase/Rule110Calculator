package main

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	"strconv"
	"time"
)

var sip [][]int

func main() {
	var data = make([]time.Duration, 10000)
	bar := progressbar.Default(10000)
	for i := 0; i < 10000; i++ {
		err := bar.Add(1)
		if err != nil {
			return
		}
		data[i] = simulate(false, 20000)
	}
	var s = time.Duration(0)
	for i := range data {
		s += data[i]
	}

	//sip = make([][]int, 20000)
	//for i := range sip {
	//	sip[i] = make([]int, 20000)
	//}
	//rand.Seed(time.Now().Unix())
	//sip[0] = rand.Perm(20000)
	//
	//bar := progressbar.Default(10000)
	//var data = make([]time.Duration, 10000)
	//for i:=0; i< 10000; i++ {
	//	err := bar.Add(1)
	//	if err != nil {
	//		return
	//	}
	//	data[i] = testoverwritespeed()
	//}
	//var s = time.Duration(0)
	//for i := range data {
	//	s += data[i]
	//}

	fmt.Println("Average time in nanoseconds: " + fmt.Sprintf("%f", float64(s.Nanoseconds())/float64(len(data))))
	fmt.Println("Average time in microseconds: " + fmt.Sprintf("%f", float64(s.Microseconds())/float64(len(data))))
	fmt.Println("Average time in milliseconds: " + fmt.Sprintf("%f", float64(s.Milliseconds())/float64(len(data))))
	fmt.Println("Average time in seconds: " + fmt.Sprintf("%f", s.Seconds()/float64(len(data))))
	fmt.Println("Total amount of simulations: " + strconv.Itoa(len(data)))

	/** SPEED TESTS (Sample: 10000)
		History TRUE @ 20000 evolutions (130 it/s)
		Average time in nanoseconds: 9626503.592100
		Average time in microseconds: 9626.503500
		Average time in milliseconds: 9.626500
		Average time in seconds: 0.009627
		Total amount of simulations: 10000
	------------------------------------------------
		History FALSE @ 20000 evolutions (83 it/s)
		Average time in nanoseconds: 11988108.519800
		Average time in microseconds: 11988.108500
		Average time in milliseconds: 11.988100
		Average time in seconds: 0.011988
		Total amount of simulations: 10000
	------------------------------------------------
		NEW! History FALSE @ 20000 evolutions (101 it/s) 21% FASTER
		Average time in nanoseconds: 9880881.055400
		Average time in microseconds: 9880.881000
		Average time in milliseconds: 9.880800
		Average time in seconds: 0.009881
		Total amount of simulations: 10000
	*/
}
