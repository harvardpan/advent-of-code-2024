package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"time"
)

func timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", name, time.Since(start))
	}
}

type Coordinate struct {
	x int
	y int
}

func main() {
	part1()
	part2()
}

func markGrid(grid [][]rune, coord Coordinate) {
	if (coord.x < 0 || coord.y < 0) || (coord.x >= len(grid[0]) || coord.y >= len(grid)) {
		// Don't mark invalid coordinates
		return
	}
	// Mark the grid with the given coordinate
	grid[coord.y][coord.x] = '#'
}

func countAntiNodes(grid [][]rune) int {
	// Count the number of anti-nodes in the grid
	result := 0
	for _, row := range grid {
		for _, char := range row {
			if char == '#' {
				result += 1
			}
		}
	}
	return result
}

func calculateAntiNodeCoordinates(coord1 Coordinate, coord2 Coordinate, rows int, columns int) (Coordinate, Coordinate) {
	// Given two coordinates, calculate the anti-node coordinates
	// The anti-node coordinates are the two coordinates that are diagonally opposite to the given coordinates
	xDiff := int(math.Abs(float64(coord1.x - coord2.x)))
	yDiff := int(math.Abs(float64(coord1.y - coord2.y)))
	antiNode1 := Coordinate{-1, -1}
	antiNode2 := Coordinate{-1, -1}
	if coord1.x < coord2.x {
		antiNode1.x = coord1.x - xDiff
		antiNode2.x = coord2.x + xDiff
	} else {
		antiNode1.x = coord1.x + xDiff
		antiNode2.x = coord2.x - xDiff
	}
	if coord1.y < coord2.y {
		antiNode1.y = coord1.y - yDiff
		antiNode2.y = coord2.y + yDiff
	} else {
		antiNode1.y = coord1.y + yDiff
		antiNode2.y = coord2.y - yDiff
	}
	// Check if the anti-node coordinates are within the grid
	if antiNode1.x < 0 || antiNode1.y < 0 || antiNode1.x >= columns || antiNode1.y >= rows {
		antiNode1 = Coordinate{-1, -1}
	}
	if antiNode2.x < 0 || antiNode2.y < 0 || antiNode2.x >= columns || antiNode2.y >= rows {
		antiNode2 = Coordinate{-1, -1}
	}
	return antiNode1, antiNode2
}

func part1() {
	// https://adventofcode.com/2024/day/7
	// Walk an operations tree and determine a valid order of operations
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	defer timer("part1")()
	result := 0
	// variables specific to this problem
	var grid [][]rune                       // stores the 2D array of characters
	antennas := make(map[rune][]Coordinate) // stores a list of coordinates for each antenna
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		var row []rune
		for colIndex, char := range line {
			row = append(row, char)
			if char != '.' {
				// Check the antennas map if it already has this key
				_, exists := antennas[char]
				if !exists {
					antennas[char] = []Coordinate{}
				}
				// Add the location to the list
				antennas[char] = append(antennas[char], Coordinate{colIndex, len(grid)})
			}
		}
		grid = append(grid, row)
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file: ", err)
	}
	// We should have the grid and the antenna locations now. Let's calculate the anti-node coordinates
	// for each pair of antennas
	for _, antenna := range antennas {
		// We need to calculate the anti-node coordinates for each pair of antennas within each type
		for i := 0; i < len(antenna); i++ {
			for j := i + 1; j < len(antenna); j++ {
				antiNode1, antiNode2 := calculateAntiNodeCoordinates(antenna[i], antenna[j], len(grid), len(grid[0]))
				markGrid(grid, antiNode1)
				markGrid(grid, antiNode2)
			}
		}
	}
	result = countAntiNodes(grid)
	fmt.Println("The final result is: ", result)
}

func calculateAllAntiNodes(coord1 Coordinate, coord2 Coordinate, rows int, columns int) []Coordinate {
	// Given two coordinates, calculate the anti-node coordinates
	// The anti-node coordinates all all the coordinates on the grip on the same slope
	xDiff := coord1.x - coord2.x
	yDiff := coord1.y - coord2.y
	antinodes := []Coordinate{coord1, coord2}
	x := coord1.x
	y := coord1.y
	for {
		x += xDiff
		y += yDiff
		if x < 0 || y < 0 || x >= columns || y >= rows {
			break
		}
		antinodes = append(antinodes, Coordinate{x, y})
	}
	for {
		x -= xDiff
		y -= yDiff
		if x < 0 || y < 0 || x >= columns || y >= rows {
			break
		}
		antinodes = append(antinodes, Coordinate{x, y})
	}
	return antinodes
}

func part2() {
	// https://adventofcode.com/2024/day/8#part2
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	defer timer("part2")()
	result := 0
	// variables specific to this problem
	var grid [][]rune                       // stores the 2D array of characters
	antennas := make(map[rune][]Coordinate) // stores a list of coordinates for each antenna
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		var row []rune
		for colIndex, char := range line {
			row = append(row, char)
			if char != '.' {
				// Check the antennas map if it already has this key
				_, exists := antennas[char]
				if !exists {
					antennas[char] = []Coordinate{}
				}
				// Add the location to the list
				antennas[char] = append(antennas[char], Coordinate{colIndex, len(grid)})
			}
		}
		grid = append(grid, row)
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file: ", err)
	}
	// We should have the grid and the antenna locations now. Let's calculate the anti-node coordinates
	// for each pair of antennas
	for _, antenna := range antennas {
		// We need to calculate the anti-node coordinates for each pair of antennas within each type
		for i := 0; i < len(antenna); i++ {
			for j := i + 1; j < len(antenna); j++ {
				antinodes := calculateAllAntiNodes(antenna[i], antenna[j], len(grid), len(grid[0]))
				for _, antinode := range antinodes {
					markGrid(grid, antinode)
				}
			}
		}
	}
	result = countAntiNodes(grid)
	fmt.Println("The final result is: ", result)
}
