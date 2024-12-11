package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", name, time.Since(start))
	}
}

type Stone struct {
	aggregates []int // counts of stones after each blink
}

func main() {
	part1()
	part2()
}

func getNextNumbers(stoneNumber int) []int {
	/*
		If the stone is engraved with the number 0, it is replaced by a stone engraved with the number 1.
		If the stone is engraved with a number that has an even number of digits, it is replaced by two stones. The left half of the digits are engraved on the new left stone, and the right half of the digits are engraved on the new right stone. (The new numbers don't keep extra leading zeroes: 1000 would become stones 10 and 0.)
		If none of the other rules apply, the stone is replaced by a new stone; the old stone's number multiplied by 2024 is engraved on the new stone.
	*/
	if stoneNumber == 0 {
		return []int{1}
	}
	// even number of digits
	stoneNumberString := strconv.Itoa(stoneNumber)
	if len(stoneNumberString)%2 == 0 {
		halfLength := len(stoneNumberString) / 2
		leftHalf, _ := strconv.Atoi(stoneNumberString[:halfLength])
		rightHalf, _ := strconv.Atoi(stoneNumberString[halfLength:])
		return []int{leftHalf, rightHalf}
	} else {
		return []int{stoneNumber * 2024}
	}
}

func blink(stones map[int]*Stone, currentStone int, stepsRemaining int) *Stone {
	// depth first search to build up the stones aggregations
	stone, exists := stones[currentStone]
	if !exists {
		// we haven't seen this stone before
		stone = &Stone{aggregates: make([]int, 0)}
		(*stone).aggregates = append((*stone).aggregates, 1) // self at 0-index
		stones[currentStone] = stone
	}
	if len((*stone).aggregates) > stepsRemaining {
		// we've already calculated this stone's aggregates. This will also be the terminating condition
		// for leaf nodes, as stepsRemaining will be 0, and each Stone will have at least one entry
		return stone
	}
	nextNumbers := getNextNumbers(currentStone)
	newAggregates := make([]int, stepsRemaining+1)
	newAggregates[0] = 1 // self at 0-index
	for _, nextNumber := range nextNumbers {
		nextStone := blink(stones, nextNumber, stepsRemaining-1)
		for i := 0; i < stepsRemaining; i++ {
			newAggregates[i+1] += (*nextStone).aggregates[i]
		}
	}
	(*stone).aggregates = newAggregates
	return stone
}

func part1() {
	// https://adventofcode.com/2024/day/11
	// Do it 25 times
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	defer timer("part1")()
	result := 0
	// variables specific to this problem
	stones := make(map[int]*Stone) // keeps track of all the stones we've "encountered"
	totalSteps := 25

	// Begin file parsing
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Day-specific code
		stoneInputs := strings.Split(line, " ")
		for _, stoneInput := range stoneInputs {

			stoneNumber, _ := strconv.Atoi(stoneInput)
			stone := blink(stones, stoneNumber, totalSteps)
			result += (*stone).aggregates[totalSteps]
		}
		// Only one line today
		break
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file: ", err)
	}
	// Post file-processing code.

	fmt.Println("The final result is: ", result)
}

func part2() {
	// https://adventofcode.com/2024/day/11#part2
	// Do it 75 times.
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	defer timer("part2")()
	result := 0
	// variables specific to this problem
	stones := make(map[int]*Stone) // keeps track of all the stones we've "encountered"
	totalSteps := 75

	// Begin file parsing
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Day-specific code
		stoneInputs := strings.Split(line, " ")
		for _, stoneInput := range stoneInputs {

			stoneNumber, _ := strconv.Atoi(stoneInput)
			stone := blink(stones, stoneNumber, totalSteps)
			result += (*stone).aggregates[totalSteps]
		}
		// Only one line today
		break
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file: ", err)
	}
	// Post file-processing code.

	fmt.Println("The final result is: ", result)
}
