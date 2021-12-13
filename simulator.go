package main

import (
	"encoding/gob"
	"fmt"
	"math"
	"os"
	"runtime"
	"strconv"
	"time"
)

var sim [][]uint64

// make sure to initialize with a raw array of 1 and 0
func initialize(initial []uint64, evolutions int, history bool) {
	for {
		if initial[0] == 0 {
			initial = initial[1:]
			continue
		} else if initial[len(initial)-1] == 0 {
			initial = initial[:len(initial)-1]
			continue
		}
		break
	}

	sim = nil
	loopLength := int(math.Ceil(float64((len(initial))+evolutions) / 64.0))

	if history == true {
		sim = make([][]uint64, evolutions+1)
		for i := range sim {
			sim[i] = make([]uint64, loopLength)
		}
	} else {
		sim = make([][]uint64, 2)
		for i := range sim {
			sim[i] = make([]uint64, loopLength)
		}
	}

	for i, j, k, p := 0, 0, 0, 0; i < 64; i++ {
		for {
			if k >= ((len(sim[0]))*64 - len(initial)) { // potential issue here with >=
				if p < len(initial) && initial[p] == 1 {
					sim[0][j] = setBit(sim[0][j], uint64(i))
				}
				p++
			}
			k++
			j++
			if j == len(sim[0]) {
				j = 0
				break
			}
		}
	}
}

func simulate(history bool, evolutions int, conditions []uint64) time.Duration {
	generationProgress = 0
	if conditions == nil {
		fmt.Println("Program uninitialized using default...")
		fmt.Println()
		initialize(r110Default(), evolutions, history)
	} else {
		initialize(conditions, evolutions, history)
	}

	start := time.Now()
	if history == true {
		historicallyAware(evolutions)
	} else {
		historicallyUnaware(evolutions)
	}
	return time.Since(start)
}

func historicallyAware(evolutions int) {
	for i := 1; i < evolutions+1; i++ {
		sim[i][0] = ((^(sim[i-1][len(sim[i-1])-1]) << 1) & sim[i-1][0]) | (sim[i-1][0] ^ sim[i-1][1])
		for j := 1; j < len(sim[i])-1; j++ {
			sim[i][j] = ((^(sim[i-1][j-1])) & sim[i-1][j]) | (sim[i-1][j] ^ sim[i-1][j+1])
		}
		sim[i][len(sim[i])-1] = ((^(sim[i-1][len(sim[i])-2])) & sim[i-1][len(sim[i])-1]) | (sim[i-1][len(sim[i])-1] ^ sim[i-1][0]>>1)
		generationProgress++
	}
}

func historicallyUnaware(evolutions int) {
	for i, k := 1, 0; i < evolutions+1; i, k = i+1, k^1 {
		sim[k^1][0] = ((^(sim[k][len(sim[k])-1]) << 1) & sim[k][0]) | (sim[k][0] ^ sim[k][1])
		for j := 1; j < len(sim[0])-1; j++ {
			sim[k^1][j] = ((^(sim[k][j-1])) & sim[k][j]) | (sim[k][j] ^ sim[k][j+1])
		}
		sim[k^1][len(sim[0])-1] = ((^(sim[k][len(sim[0])-2])) & sim[k][len(sim[0])-1]) | (sim[k][len(sim[0])-1] ^ sim[k][0]>>1)
		generationProgress++
	}
}

func decompress(history bool) [][]uint64 {
	// TODO reformat these loops in compliance with initialization function
	if history {
		out := make([][]uint64, len(sim))
		for i := range out {
			out[i] = make([]uint64, len(sim[0])*64)
		}
		lpc := 0
		for i := 0; i < len(sim); i++ {
			for j := 0; j < 64; j++ {
				for k := 0; k < len(sim[i]); k++ {
					out[i][lpc] = uint64(getBit(sim[i][k], j))
					lpc++
				}
			}
			lpc = 0
		}
		return out
	} else {
		out := make([][]uint64, 1)
		for i := range out {
			out[i] = make([]uint64, len(sim[0])*64)
		}
		if countLeadingZeros(0) < countLeadingZeros(1) {
			lpc := 0
			for j := 0; j < 64; j++ {
				for k := 0; k < len(sim[0]); k++ {
					out[0][lpc] = uint64(getBit(sim[0][k], j))
					lpc++
				}
			}
		} else {
			lpc := 0
			for j := 0; j < 64; j++ {
				for k := 0; k < len(sim[1]); k++ {
					out[1][lpc] = uint64(getBit(sim[1][k], j))
					lpc++
				}
			}
		}
		return out
	}
}

func countLeadingZeros(layer int) int {
	zeros := 0
	for j := 0; j < 64; j++ {
		for k := 0; k < len(sim[layer]); k++ {
			if getBit(sim[layer][k], j) == 1 {
				return zeros
			}
			zeros++
		}
	}
	return -1
}

func readTape() {
}

func writeToFile(path string, obj interface{}) {
	file, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	enc := gob.NewEncoder(file)
	if err = enc.Encode(&obj); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Structure written into file successfully")

	err = file.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func readFromFile(path string, assign interface{}) interface{} {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	dec := gob.NewDecoder(file)
	if err = dec.Decode(&assign); err != nil {
		fmt.Println(err)
	}
	fmt.Println("Structure read from file successfully")
	return assign
}

// REMEMBER: BINARY NUMBERS READ RIGHT TO LEFT!!!
func displayRaw(layer int) {
	for i := 0; i < len(sim[layer]); i++ {
		fmt.Printf("%064d", strconv.FormatInt(int64(sim[layer][i]), 2))
		fmt.Println()
	}
}

func displayFancy() {
	if sim == nil {
		fmt.Println("Array Empty!")
	}
	for i := 0; i < len(sim); i++ {
		for j := 0; j < 64; j++ {
			for k := 0; k < len(sim[i]); k++ {
				if getBit(sim[i][k], j) == 0 {
					fmt.Print("□")
				} else {
					fmt.Print("■")
				}
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func setBit(n uint64, pos uint64) uint64 {
	n |= 1 << pos
	return n
}

func clearBit(n uint64, pos uint64) uint64 {
	mask := ^(1 << pos)
	n &= uint64(mask)
	return n
}

func getBit(n uint64, pos int) int {
	val := n & (1 << pos)
	if val > 0 {
		return 1
	}
	return 0
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
