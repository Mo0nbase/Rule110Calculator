package main

import (
	"fmt"
	"math"
	"strconv"
)

var sim [][]uint64

/*
	NOTE: presently this function preallocates all needed ram for the simulation this could be a potential problem if the user is facing memory constraints.
	A potential fix is to simply detect if the current evolution is divisible by 64 and within the copying phase make sure that the program starts the evolution
	from w[1] instead of w[0] because in all cases with a rightmost set bit the simulation must expand by 1. This would also dynamically allocate memory as the program
	progresses however would require an additional comparison and overflow bit if it turns out that the first bit of the getBit(sim[i-1][0],0)==1. But although this
	would solve the problem for the first bit if we simply append an empty uint64 with 64 0's each of those zeros will be used as the leftmost neighbour of w[i][1] which
	completely borks the simulation. Not only must we append an extra bit to the beginning of the theoretical w[i][0] we must also make sure it doesn't interfere
	with the other bits neighbours. A method for this is not presently known. Further, research required.
		** Also, error prediction must be built in to detect if the predicted amount of ram pessary for the simulation can be met by the system.

	NOTE: presently this function for every 64 evolutions adds a 64-bit integer to make more room for boundary overflow however the simulation function
	is built to begin each evolution on the leftmost column and because rule 110 only grows to the left this result in a waste of ram in addition to violating
	the boundary conditions for the rest of the simulation. A potential fix is to detect if the current evolution number is divisible by 64 and starting the simulation
	on w[1] instead of w[0]. But this problem might also resolve itself due to the order in which bits are processed from right to left. Further, testing required.
	This solution would also interfere with a problem discussed later in the first note above regarding the neighbours of bits on a row > 0.
*/
// make sure to initialize with a raw array of 1 and 0
func initialize(initial []uint64, evolutions int) {
	loopLength := int(math.Ceil(float64(len(initial) / 64)))

	for i := 0; i <= evolutions; i++ {
		sim = append(sim, []uint64{})
	}

	if loopLength == 0 {
		sim[0] = append(sim[0], uint64(0))
		for i := 1; i < evolutions; i++ {
			for j := 0; j < loopLength+(i+1/64); j++ {
				sim[i] = append(sim[i], uint64(0))
			}
		}
		for i := 0; i < 64; i++ {
			if i < len(initial) && initial[i] == 1 {
				sim[0][0] = setBit(sim[0][0], uint64(i))
			}
		}
	} else {
		for i := 0; i < loopLength; i++ {
			sim[0] = append(sim[0], uint64(0))
		}
		for i := 1; i < evolutions; i++ {
			for j := 0; j < loopLength+(i+1/64); j++ {
				sim[i] = append(sim[i], uint64(0))
			}
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
		initialize(testSplit, evolutions)
	}

	if history == true {
		for i := 1; i < evolutions+1; i++ {
			for j := 0; j < len(sim[i]); j++ {
				// NOTE: This if statement handles the starting ending and general boundary conditions for the program
				// however there is potential confusion going on with the shift operations. Due to the order integer's bits are read
				// by the CPU there is a potential issue going on with the direction (I.E. <<,>>) of the shifts being used
				// from the opposite side of the present index of the array (0 or n-1). If the simulation is not displaying correctly
				// the first fix that can be tried is just reversing the direction of all shift operators.
				if j == 0 {
					sim[i][j] = (sim[i-1][len(sim[i-1])]>>1 | sim[i-1][j]) ^ (sim[i-1][len(sim[i-1])] >> 1 & sim[i-1][j] & sim[i-1][j+1])
				} else if j == len(sim[i])-1 {
					sim[i][j] = (sim[i-1][j-1] | sim[i-1][j]) ^ (sim[i-1][j-1] & sim[i-1][j] & sim[i-1][0] << 1)
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
	//displayFancy()
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
