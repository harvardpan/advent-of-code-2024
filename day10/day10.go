package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
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

type Location struct {
	x      int
	y      int
	height int
}

func safeGet(grid [][]Location, row, col int) Location {
	if row < 0 || col < 0 {
		return Location{x: -1, y: -1, height: -1}
	}
	if row >= len(grid) || col >= len(grid[0]) {
		return Location{x: -1, y: -1, height: -1}
	}
	return grid[row][col]
}

func dfs(grid [][]Location, x, y int, visited map[string]bool, result *int) {
	if x < 0 || y < 0 || x >= len(grid) || y >= len(grid[0]) {
		return
	}
	if visited != nil {
		if visited[fmt.Sprintf("%d,%d", x, y)] {
			return
		}
		visited[fmt.Sprintf("%d,%d", x, y)] = true
	}

	if grid[x][y].height == 9 {
		*result++
		return
	}

	currentHeight := grid[x][y].height
	if safeGet(grid, x-1, y).height == currentHeight+1 {
		dfs(grid, x-1, y, visited, result)
	}
	if safeGet(grid, x+1, y).height == currentHeight+1 {
		dfs(grid, x+1, y, visited, result)
	}
	if safeGet(grid, x, y-1).height == currentHeight+1 {
		dfs(grid, x, y-1, visited, result)
	}
	if safeGet(grid, x, y+1).height == currentHeight+1 {
		dfs(grid, x, y+1, visited, result)
	}
}

func part1() {
	// https://adventofcode.com/2024/day/10
	// Find the number of trailhead to peaks
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	defer timer("part1")()
	result := 0
	// variables specific to this problem
	var grid [][]Location // stores the 2D array of locations
	// Begin file parsing
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		rowIndex := 0
		line := scanner.Text()
		var row []Location
		for colIndex, char := range line {
			// convert char to integer
			height, _ := strconv.Atoi(string(char))
			row = append(row, Location{x: colIndex, y: rowIndex, height: height})
		}
		grid = append(grid, row)
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file: ", err)
	}
	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[0]); j++ {
			visited := make(map[string]bool)
			if grid[i][j].height == 0 {
				dfs(grid, i, j, visited, &result)
			}
		}
	}
	fmt.Println("The final result is: ", result)
}

func part2() {
	// https://adventofcode.com/2024/day/10#part2
	// Find number of distinct paths to the same destination
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	defer timer("part1")()
	result := 0
	// variables specific to this problem
	var grid [][]Location // stores the 2D array of locations
	// Begin file parsing
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		rowIndex := 0
		line := scanner.Text()
		var row []Location
		for colIndex, char := range line {
			// convert char to integer
			height, _ := strconv.Atoi(string(char))
			row = append(row, Location{x: colIndex, y: rowIndex, height: height})
		}
		grid = append(grid, row)
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file: ", err)
	}
	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[0]); j++ {
			if grid[i][j].height == 0 {
				dfs(grid, i, j, nil, &result)
			}
		}
	}
	fmt.Println("The final result is: ", result)
}
