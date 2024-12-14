package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"time"

	"gonum.org/v1/gonum/mat"
)

func timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", name, time.Since(start))
	}
}

func main() {
	part1()
	part2()
}

func isWholeNumber(f float64) bool {
	return f == float64(int(f))
}

func solveEquations(x1, y1, c1, x2, y2, c2 float64) (int, int) {
	// Solve for A and B in the following equations:
	// x1 * A + y1 * A = c1
	// x2 * B + y2 * B = c2
	fmt.Println("Solving for A and B in the following equations:")
	fmt.Println(x1, " * A + ", y1, " * A = ", c1)
	fmt.Println(x2, " * B + ", y2, " * B = ", c2)
	// Create A matrix and solve for x and y
	A := mat.NewDense(2, 2, []float64{x1, x2, y1, y2}) // co-efficient matrix
	b := mat.NewVecDense(2, []float64{c1, c2})         // constant vector
	var X mat.VecDense
	X.SolveVec(A, b) // Solving for A * X = b
	// Even though we solve for floating point numbers, button presses are integers.
	roundedAPresses := math.Round(X.At(0, 0)*1000) / 1000
	if !isWholeNumber(roundedAPresses) {
		roundedAPresses = 0
	}
	roundedBPresses := math.Round(X.At(1, 0)*1000) / 1000
	if !isWholeNumber(roundedBPresses) {
		roundedBPresses = 0
	}
	if roundedAPresses <= 0 || roundedBPresses <= 0 {
		fmt.Println("No solution found for this prize.")
	} else {
		fmt.Println("A: ", X.At(0, 0), " B: ", X.At(1, 0))
	}
	return int(roundedAPresses), int(roundedBPresses)
}

func part1() {
	// https://adventofcode.com/2024/day/13
	// Do matrix math to solve two linear equations
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	defer timer("part1")()
	result := 0
	// variables specific to this problem
	buttonPattern := regexp.MustCompile(`Button (\w): X\+(\d+), Y\+(\d+)`)
	prizePattern := regexp.MustCompile(`Prize: X=(\d+), Y=(\d+)`)
	buttons := make([][2]float64, 0)

	// Begin file parsing
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Day-specific code
		if buttonPattern.MatchString(line) {
			matches := buttonPattern.FindStringSubmatch(line)
			x, err := strconv.ParseFloat(matches[2], 64)
			if err != nil {
				fmt.Println("Error converting X portion:", err)
				return
			}
			y, err := strconv.ParseFloat(matches[3], 64)
			if err != nil {
				fmt.Println("Error converting Y portion:", err)
				return
			}
			buttons = append(buttons, [2]float64{x, y})
		} else if prizePattern.MatchString(line) {
			matches := prizePattern.FindStringSubmatch(line)
			c1, err := strconv.ParseFloat(matches[1], 64)
			if err != nil {
				fmt.Println("Error converting X prize:", err)
				return
			}
			c2, err := strconv.ParseFloat(matches[2], 64)
			if err != nil {
				fmt.Println("Error converting Y prize:", err)
				return
			}
			// Do something with the prize coordinates
			aPresses, bPresses := solveEquations(buttons[0][0], buttons[0][1], c1, buttons[1][0], buttons[1][1], c2)
			// Reset buttons
			buttons = make([][2]float64, 0)

			if aPresses <= 0 || aPresses > 100 || bPresses <= 0 || bPresses > 100 {
				// No solution, simply continue
				continue
			}
			fmt.Println("A presses: ", aPresses, " B presses: ", bPresses)
			result += aPresses*3 + bPresses
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file: ", err)
	}
	// Post file-processing code.
	fmt.Println("The final result is: ", result)
}

func part2() {
	// https://adventofcode.com/2024/day/13#part2
	// Do matrix math to solve two linear equations, with slight modification to conditions
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	defer timer("part2")()
	result := 0
	// variables specific to this problem
	buttonPattern := regexp.MustCompile(`Button (\w): X\+(\d+), Y\+(\d+)`)
	prizePattern := regexp.MustCompile(`Prize: X=(\d+), Y=(\d+)`)
	buttons := make([][2]float64, 0)

	// Begin file parsing
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Day-specific code
		if buttonPattern.MatchString(line) {
			matches := buttonPattern.FindStringSubmatch(line)
			x, err := strconv.ParseFloat(matches[2], 64)
			if err != nil {
				fmt.Println("Error converting X portion:", err)
				return
			}
			y, err := strconv.ParseFloat(matches[3], 64)
			if err != nil {
				fmt.Println("Error converting Y portion:", err)
				return
			}
			buttons = append(buttons, [2]float64{x, y})
		} else if prizePattern.MatchString(line) {
			matches := prizePattern.FindStringSubmatch(line)
			c1, err := strconv.ParseFloat(matches[1], 64)
			if err != nil {
				fmt.Println("Error converting X prize:", err)
				return
			}
			c2, err := strconv.ParseFloat(matches[2], 64)
			if err != nil {
				fmt.Println("Error converting Y prize:", err)
				return
			}
			// Part 2 increases the c1 and c2 values by 10000000000000
			c1 += 10000000000000.0
			c2 += 10000000000000.0
			// Do something with the prize coordinates
			aPresses, bPresses := solveEquations(buttons[0][0], buttons[0][1], c1, buttons[1][0], buttons[1][1], c2)
			// Reset buttons
			buttons = make([][2]float64, 0)

			if aPresses <= 0 || bPresses <= 0 {
				// No solution, simply continue
				continue
			}
			fmt.Println("A presses: ", aPresses, " B presses: ", bPresses)
			result += aPresses*3 + bPresses
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file: ", err)
	}
	// Post file-processing code.
	fmt.Println("The final result is: ", result)
}
