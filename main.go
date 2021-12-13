package main

import (
	"fmt"
	"github.com/faiface/pixel/pixelgl"
	"github.com/schollz/progressbar/v3"
	"strconv"
	"time"
)

func main() {
	pixelgl.Run(run)
	//performanceTest(4, 1000000, false)
	//readFromFile()
	//displayFancy()
	//simulate(true, 3500, randSt(3500)) // 2 below max // 7000
	//fmt.Println(calculateArea(3500, 3500, 2))
}

func calculateArea(length int, evolutions int, pSize int) int {
	return calculateWidth(length, evolutions, pSize) * calculateHeight(evolutions, pSize)
}

func calculateWidth(length int, evolutions int, pSize int) int {
	base := length + evolutions + 1
	return (base + (64 - (base % 64))) * pSize
}

func calculateHeight(evolutions int, pSize int) int {
	return (evolutions + 1) * pSize
}

func performanceTest(repetitions int, evolutions int, history bool) {
	var data = make([]time.Duration, repetitions)
	bar := progressbar.Default(int64(repetitions))
	for i := 0; i < repetitions; i++ {
		data[i] = simulate(history, evolutions, r110Default())
		err := bar.Add(1)
		if err != nil {
			return
		}
	}
	var s = time.Duration(0)
	for i := range data {
		s += data[i]
	}

	fmt.Println("Average time in nanoseconds: " + fmt.Sprintf("%f", float64(s.Nanoseconds())/float64(len(data))))
	fmt.Println("Average time in microseconds: " + fmt.Sprintf("%f", float64(s.Microseconds())/float64(len(data))))
	fmt.Println("Average time in milliseconds: " + fmt.Sprintf("%f", float64(s.Milliseconds())/float64(len(data))))
	fmt.Println("Average time in seconds: " + fmt.Sprintf("%f", s.Seconds()/float64(len(data))))
	fmt.Println("Total amount of simulations: " + strconv.Itoa(len(data)))
}

/** SPEED TESTS (Sample: 10000)
	History TRUE @ 20000 evolutions (63 it/s)
	Average time in nanoseconds: 7616993.612300
	Average time in microseconds: 7616.993600
	Average time in milliseconds: 7.616900
	Average time in seconds: 0.007617
	Total amount of simulations: 10000
------------------------------------------------
	History FALSE @ 20000 evolutions (117 it/s)
	Average time in nanoseconds: 8520770.050700
	Average time in microseconds: 8520.770000
	Average time in milliseconds: 8.520700
	Average time in seconds: 0.008521
	Total amount of simulations: 10000
*/
