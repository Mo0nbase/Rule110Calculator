package main

import (
	"fmt"
	"math"
	"strconv"
)

var sim [][]uint64

// make sure to initialize with a raw array of 1 and 0
func initialize(initial []uint64, evolutions int) {
	loopLength := int(math.Ceil(float64(len(initial) / 64)))

	for i := 0; i <= evolutions; i++ {
		sim = append(sim, []uint64{})
	}

	if loopLength == 0 {
		sim[0] = append(sim[0], uint64(0))
		for i := 0; i < 64; i++ {
			if i < len(initial) && initial[i] == 1 {
				sim[0][0] = setBit(sim[0][0], uint64(i))
			}
		}
	} else {
		for i := 0; i < loopLength; i++ {
			sim[0] = append(sim[0], uint64(0))
		}
		pos := 0
		for i := 0; i < 64; i++ {
			for j := 0; j < loopLength; j++ {
				if pos < len(initial) && initial[pos] == 1 {
					sim[0][j] = setBit(sim[0][j], uint64(i))
				}
				pos++
			}
		}
	}
}

func simulate(history bool, evolutions int) {
	if sim == nil {
		fmt.Println("Program uninitialized using default...")
		fmt.Println()
		var testSplit = r110Default()
		initialize(testSplit, 0)
	}

	if history == true {
		for i := 1; i < evolutions; i++ {
			for j := 0; j < len(sim[i]); j++ {
				if j == 0 {

				} else if j == len(sim[i])-1 {

				} else {
					// Xor[Or[p, q], And[p, q, r]]
					// Xor[p, Or[q, r]]
					// w[i] = w[i - 1] ^ (w[i] | w[i + 1])
					sim[i][j] = (sim[i-1][j-1] | sim[i-1][j]) ^ (sim[i-1][j-1] & sim[i-1][j] & sim[i-1][j+1])
				}
			}
		}
	} else {

	}
	displayFancy()
	// if history is true simulate and store entire history
}

func readTape() uint64 {
	return 1
}

func writeToFile() {

}

func readFromFile() {

}

// REMEMBER: BINARY NUMBERS READ RIGHT TO LEFT!!!
func displayRaw(layer int) {
	fmt.Println("REMEMBER: BINARY NUMBERS READ RIGHT TO LEFT!!!")
	for i := 0; i < len(sim[layer]); i++ {
		fmt.Printf("%064d", strconv.FormatInt(int64(sim[0][i]), 2))
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
				fmt.Print(" " + strconv.Itoa(getBit(sim[i][k], uint64(j))))
			}
		}
	}
}

// Sets the bit at pos in the integer n.
func setBit(n uint64, pos uint64) uint64 {
	n |= 1 << pos
	return n
}

// Clears the bit at pos in n.
func clearBit(n uint64, pos uint64) uint64 {
	mask := ^(1 << pos)
	n &= uint64(mask)
	return n
}

func getBit(n uint64, pos uint64) int {
	val := n & (1 << pos)
	if val > 0 {
		return 1
	}
	return 0
}
