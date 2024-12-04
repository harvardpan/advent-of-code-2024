package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	part1()
	part2()
}

func safeGet(grid [][]rune, row, col int) rune {
	if row < 0 || col < 0 {
		return ' '
	}
	if row >= len(grid) || col >= len(grid[0]) {
		return ' '
	}
	return grid[row][col]
}

func calculateXmasAllDirections(grid [][]rune) int {
	result := 0
	numrows := len(grid)
	numcols := len(grid[0])
	for row := 0; row < numrows; row++ {
		for col := 0; col < numcols; col++ {
			if grid[row][col] != 'X' {
				continue
			}
			// If character is "X", then check for number of XMAS
			forward := []rune{'X', safeGet(grid, row, col+1), safeGet(grid, row, col+2), safeGet(grid, row, col+3)}
			backward := []rune{'X', safeGet(grid, row, col-1), safeGet(grid, row, col-2), safeGet(grid, row, col-3)}
			up := []rune{'X', safeGet(grid, row-1, col), safeGet(grid, row-2, col), safeGet(grid, row-3, col)}
			down := []rune{'X', safeGet(grid, row+1, col), safeGet(grid, row+2, col), safeGet(grid, row+3, col)}
			se_diag := []rune{'X', safeGet(grid, row+1, col+1), safeGet(grid, row+2, col+2), safeGet(grid, row+3, col+3)}
			sw_diag := []rune{'X', safeGet(grid, row+1, col-1), safeGet(grid, row+2, col-2), safeGet(grid, row+3, col-3)}
			ne_diag := []rune{'X', safeGet(grid, row-1, col+1), safeGet(grid, row-2, col+2), safeGet(grid, row-3, col+3)}
			nw_diag := []rune{'X', safeGet(grid, row-1, col-1), safeGet(grid, row-2, col-2), safeGet(grid, row-3, col-3)}
			for _, direction := range [][]rune{forward, backward, up, down, se_diag, sw_diag, ne_diag, nw_diag} {
				if string(direction) == "XMAS" {
					result += 1
				}
			}
		}
	}
	return result
}

func calculateMasInXFormation(grid [][]rune) int {
	result := 0
	numrows := len(grid)
	numcols := len(grid[0])
	for row := 0; row < numrows; row++ {
		for col := 0; col < numcols; col++ {
			if grid[row][col] != 'A' {
				continue
			}
			// If character is "A", then check for two 3-letter diagonals
			ne_sw_diag := []rune{safeGet(grid, row-1, col+1), 'A', safeGet(grid, row+1, col-1)}
			sw_ne_diag := []rune{safeGet(grid, row+1, col-1), 'A', safeGet(grid, row-1, col+1)}
			nw_se_diag := []rune{safeGet(grid, row-1, col-1), 'A', safeGet(grid, row+1, col+1)}
			se_nw_diag := []rune{safeGet(grid, row+1, col+1), 'A', safeGet(grid, row-1, col-1)}
			diag_count := 0
			if string(ne_sw_diag) == "MAS" || string(sw_ne_diag) == "MAS" {
				diag_count += 1
			}
			if string(nw_se_diag) == "MAS" || string(se_nw_diag) == "MAS" {
				diag_count += 1
			}
			if diag_count == 2 {
				result += 1
			}
		}
	}
	return result
}

func part1() {
	// https://adventofcode.com/2024/day/4
	// Do a word search of all XMAS in the grid
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	var grid [][]rune // stores the 2D array of characters
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		var row []rune
		for _, char := range line {
			row = append(row, char)
		}
		grid = append(grid, row)
	}
	result := calculateXmasAllDirections(grid)
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file: ", err)
	}
	fmt.Println("The final result is: ", result)
}

func part2() {
	// https://adventofcode.com/2024/day/4#part2
	// Find all MAS in a cross pattern
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	var grid [][]rune // stores the 2D array of characters
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		var row []rune
		for _, char := range line {
			row = append(row, char)
		}
		grid = append(grid, row)
	}
	result := calculateMasInXFormation(grid)
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file: ", err)
	}
	fmt.Println("The final result is: ", result)
}
