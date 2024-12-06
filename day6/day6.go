package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"
)

func main() {
	part1()
	part2()
}

func deepCopyGrid(grid [][]rune) [][]rune {
	// Create a new grid with the same dimensions
	newGrid := make([][]rune, len(grid))
	for i := range grid {
		newGrid[i] = make([]rune, len(grid[i]))
		copy(newGrid[i], grid[i])
	}
	return newGrid
}

func walkGrid(grid [][]rune, posRow, posCol int, direction rune) int {
	visited := make(map[string]bool)
	visited[fmt.Sprintf("%d,%d,%c", posRow, posCol, direction)] = true
	result := 1
	grid[posRow][posCol] = 'X' // mark the initial position as visited
	for {
		nextRow := -1
		nextCol := -1
		// Based on the direction, we determine what the next position is supposed to be.
		switch direction {
		case 'n':
			nextRow = posRow - 1
			nextCol = posCol
		case 'e':
			nextRow = posRow
			nextCol = posCol + 1
		case 's':
			nextRow = posRow + 1
			nextCol = posCol
		case 'w':
			nextRow = posRow
			nextCol = posCol - 1
		}
		if nextRow < 0 || nextRow >= len(grid) || nextCol < 0 || nextCol >= len(grid[0]) {
			// Ending condition - we've passed the boundary of the grid.
			return result
		}
		if grid[nextRow][nextCol] == '#' {
			// We've hit a barrier. Change direction
			switch direction {
			case 'n':
				direction = 'e'
			case 'e':
				direction = 's'
			case 's':
				direction = 'w'
			case 'w':
				direction = 'n'
			}
		} else {
			// No barrier, can continue
			posRow = nextRow
			posCol = nextCol
			// Check if we've visited this position before in the same direction. If so,
			// we're in an infinite loop. Return -1 to indicate that.
			if visited[fmt.Sprintf("%d,%d,%c", posRow, posCol, direction)] {
				return -1
			}
			visited[fmt.Sprintf("%d,%d,%c", posRow, posCol, direction)] = true
			if grid[posRow][posCol] == '.' {
				grid[posRow][posCol] = 'X'
				result++
			}
		}
	}
}

func part1() {
	// https://adventofcode.com/2024/day/6
	// Find path through the grid while navigating barriers
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	result := 0
	var grid [][]rune // stores the 2D array of characters
	posRow := -1
	posCol := -1
	rowIndex := 0
	direction := 'n' // n, e, s, w are the options
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		var row []rune
		for colIndex, char := range line {
			row = append(row, char)
			if char == '^' {
				posRow = rowIndex
				posCol = colIndex
			}
			colIndex++
		}
		grid = append(grid, row)
		rowIndex++
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file: ", err)
	}
	defer timer("part1")()
	// We should have the grid and the starting position now. Let's navigate the grid
	gridCopy := deepCopyGrid(grid)
	result = walkGrid(gridCopy, posRow, posCol, direction)
	fmt.Println("The final result is: ", result)
}

func timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", name, time.Since(start))
	}
}

func part2() {
	// https://adventofcode.com/2024/day/6#part2
	// Loop through and find the number of places we can place a barrier to get an infinite loop
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	result := 0
	var grid [][]rune // stores the 2D array of characters
	posRow := -1
	posCol := -1
	rowIndex := 0
	direction := 'n' // n, e, s, w are the options
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		var row []rune
		for colIndex, char := range line {
			row = append(row, char)
			if char == '^' {
				posRow = rowIndex
				posCol = colIndex
			}
			colIndex++
		}
		grid = append(grid, row)
		rowIndex++
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file: ", err)
	}
	defer timer("part2")()
	// Add goroutine to speed up the calculation
	var wg sync.WaitGroup
	var mu sync.Mutex
	// We should have the grid and the starting position now. Let's navigate the grid
	for rowIndex, row := range grid {
		for colIndex, char := range row {
			if char == '.' {
				wg.Add(1) // add to the wait group
				go func(rowIndex, colIndex int) {
					defer wg.Done() // defer the done call
					gridCopy := deepCopyGrid(grid)
					gridCopy[rowIndex][colIndex] = '#'
					if walkGrid(gridCopy, posRow, posCol, direction) == -1 {
						mu.Lock() // prevent concurrent writes to result
						result++
						mu.Unlock()
						// fmt.Println("Placing a barrier at (", rowIndex, ",", colIndex, ") will create an infinite loop.")
					}
				}(rowIndex, colIndex)
			}
		}
	}
	wg.Wait() // wait for all goroutines to finish
	fmt.Println("The final result is: ", result)
}
