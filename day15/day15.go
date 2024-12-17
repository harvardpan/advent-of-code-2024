package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"time"
)

func timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", name, time.Since(start))
	}
}

type Position struct {
	x       int
	y       int
	content rune
}

func main() {
	part1()
	part2()
}

func printGrid(grid [][]*Position) {
	for _, row := range grid {
		for _, position := range row {
			fmt.Print(string(position.content))
		}
		fmt.Println()
	}
}

func getNextPosition(grid [][]*Position, currentPosition Position, instruction rune) Position {
	nextPosition := Position{x: currentPosition.x, y: currentPosition.y}
	switch instruction {
	case '>':
		nextPosition.x = currentPosition.x + 1
	case '<':
		nextPosition.x = currentPosition.x - 1
	case '^':
		nextPosition.y = currentPosition.y - 1
	case 'v':
		nextPosition.y = currentPosition.y + 1
	}
	if nextPosition.x < 0 || nextPosition.x >= len(grid[0]) || nextPosition.y < 0 || nextPosition.y >= len(grid) {
		return Position{x: -1, y: -1}
	}
	return nextPosition
}

func shiftBox(grid [][]*Position, currentPosition Position, instruction rune) bool {
	nextPosition := getNextPosition(grid, currentPosition, instruction)
	if nextPosition.x == -1 || nextPosition.y == -1 {
		return false
	}
	if grid[nextPosition.y][nextPosition.x].content == '#' {
		// Can't move an obstacle
		return false
	}
	if grid[nextPosition.y][nextPosition.x].content == '.' {
		// Found an empty space. Move the box here and return true. Terminating condition for recursion.
		grid[currentPosition.y][currentPosition.x].content = '.'
		grid[nextPosition.y][nextPosition.x].content = 'O'
		return true
	}
	if grid[nextPosition.y][nextPosition.x].content == 'O' {
		// Box at the current position. Recursively check for an empty space in the direction of the instruction
		if !shiftBox(grid, nextPosition, instruction) {
			return false
		}
		grid[currentPosition.y][currentPosition.x].content = '.'
		grid[nextPosition.y][nextPosition.x].content = 'O'
		return true
	}
	return false
}

func processInstruction(grid [][]*Position, robotPosition *Position, instruction rune) {
	nextPosition := getNextPosition(grid, *robotPosition, instruction)
	if nextPosition.x == -1 || nextPosition.y == -1 {
		return
	}
	if grid[nextPosition.y][nextPosition.x].content == '#' {
		// Can't move an obstacle
		return
	}
	if grid[nextPosition.y][nextPosition.x].content == '.' {
		grid[robotPosition.y][robotPosition.x].content = '.'
		grid[nextPosition.y][nextPosition.x].content = '@'
		robotPosition.x = nextPosition.x
		robotPosition.y = nextPosition.y
		return
	}
	if grid[nextPosition.y][nextPosition.x].content == 'O' {
		// Can move into this spot as long as we can shift all the boxes in the direction of the instruction
		if !shiftBox(grid, nextPosition, instruction) {
			return
		}
		grid[robotPosition.y][robotPosition.x].content = '.'
		grid[nextPosition.y][nextPosition.x].content = '@'
		robotPosition.x = nextPosition.x
		robotPosition.y = nextPosition.y
	}
}

func calculateScore(grid [][]*Position) int {
	result := 0
	for _, row := range grid {
		for _, position := range row {
			if position.content == 'O' {
				result += position.x + 100*position.y
			}
		}
	}
	return result
}

func part1() {
	// https://adventofcode.com/2024/day/15
	// Move boxes around using a robot
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	defer timer("part1")()
	result := 0
	// variables specific to this problem
	topBottomPattern := regexp.MustCompile(`^#+$`)
	gridRowPattern := regexp.MustCompile(`^#[.O@#]+#$`)
	instructionPattern := regexp.MustCompile(`[<v>^]+`)
	var grid [][]*Position
	var robotPosition Position
	fullInstructions := ""

	// Begin file parsing
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Day-specific code
		row := make([]*Position, 0)
		if topBottomPattern.MatchString(line) || gridRowPattern.MatchString(line) {
			for colIndex, char := range line {
				position := Position{x: colIndex, y: len(grid), content: char}
				row = append(row, &position)
				if string(char) == "@" {
					robotPosition = position // make a copy
				}
			}
			grid = append(grid, row)
		} else if instructionPattern.MatchString(line) {
			fullInstructions += line
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file: ", err)
	}
	// Post file-processing code.
	printGrid(grid)
	// Go through each instruction
	for _, instruction := range fullInstructions {
		processInstruction(grid, &robotPosition, instruction)
	}
	result = calculateScore(grid)
	fmt.Println("The final result is: ", result)
}

func processInstruction2(grid [][]*Position, robotPosition *Position, instruction rune) {
	nextPosition := getNextPosition(grid, *robotPosition, instruction)
	if nextPosition.x == -1 || nextPosition.y == -1 {
		return
	}
	if grid[nextPosition.y][nextPosition.x].content == '#' {
		// Can't move an obstacle
		return
	}
	if grid[nextPosition.y][nextPosition.x].content == '.' {
		grid[robotPosition.y][robotPosition.x].content = '.'
		grid[nextPosition.y][nextPosition.x].content = '@'
		robotPosition.x = nextPosition.x
		robotPosition.y = nextPosition.y
		return
	}
	if grid[nextPosition.y][nextPosition.x].content == '[' || grid[nextPosition.y][nextPosition.x].content == ']' {
		// Can move into this spot as long as we can shift all the boxes in the direction of the instruction
		if !areAllLeafPositionsOpen(grid, nextPosition, instruction, make(map[*Position]bool)) {
			return
		}
		shiftBox2(grid, nextPosition, instruction, make(map[*Position]bool))
		grid[robotPosition.y][robotPosition.x].content = '.'
		grid[nextPosition.y][nextPosition.x].content = '@'
		robotPosition.x = nextPosition.x
		robotPosition.y = nextPosition.y
	}
}

func areAllLeafPositionsOpen(grid [][]*Position, boxPosition Position, instruction rune, visited map[*Position]bool) bool {
	if visited[grid[boxPosition.y][boxPosition.x]] {
		// already visited this position, so assume true.
		return true
	}
	visited[grid[boxPosition.y][boxPosition.x]] = true
	// Check leaf condition that would terminate the recursion
	if grid[boxPosition.y][boxPosition.x].content == '#' {
		return false
	} else if grid[boxPosition.y][boxPosition.x].content == '.' {
		return true
	}
	nextPosition := getNextPosition(grid, boxPosition, instruction)

	if grid[boxPosition.y][boxPosition.x].content == '[' {
		return areAllLeafPositionsOpen(grid, nextPosition, instruction, visited) && areAllLeafPositionsOpen(grid, Position{x: boxPosition.x + 1, y: boxPosition.y}, instruction, visited)
	} else if grid[boxPosition.y][boxPosition.x].content == ']' {
		return areAllLeafPositionsOpen(grid, nextPosition, instruction, visited) && areAllLeafPositionsOpen(grid, Position{x: boxPosition.x - 1, y: boxPosition.y}, instruction, visited)
	}
	return false
}

func shiftBox2(grid [][]*Position, boxPosition Position, instruction rune, visited map[*Position]bool) {
	if visited[grid[boxPosition.y][boxPosition.x]] {
		// already visited this position, so assume true.
		return
	}
	visited[grid[boxPosition.y][boxPosition.x]] = true
	if grid[boxPosition.y][boxPosition.x].content == '[' {
		shiftBox2(grid, Position{x: boxPosition.x + 1, y: boxPosition.y}, instruction, visited)
	} else if grid[boxPosition.y][boxPosition.x].content == ']' {
		shiftBox2(grid, Position{x: boxPosition.x - 1, y: boxPosition.y}, instruction, visited)
	}
	nextPosition := getNextPosition(grid, boxPosition, instruction)
	if nextPosition.x == -1 || nextPosition.y == -1 {
		return
	}
	if grid[nextPosition.y][nextPosition.x].content == '.' {
		// terminating position.
		grid[nextPosition.y][nextPosition.x].content = grid[boxPosition.y][boxPosition.x].content
		grid[boxPosition.y][boxPosition.x].content = '.'
		return
	}
	shiftBox2(grid, nextPosition, instruction, visited)
	if grid[nextPosition.y][nextPosition.x].content == '.' {
		// terminating position.
		grid[nextPosition.y][nextPosition.x].content = grid[boxPosition.y][boxPosition.x].content
		grid[boxPosition.y][boxPosition.x].content = '.'
		return
	}
}

func calculateScore2(grid [][]*Position) int {
	result := 0
	for _, row := range grid {
		for _, position := range row {
			if position.content == '[' {
				result += position.x + 100*position.y
			}
		}
	}
	return result
}

func part2() {
	// https://adventofcode.com/2024/day/15#part2
	// Move boxes around using a robot on a grid where everything is twice as wide
	// Have to account for moving multiple boxes at once
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	defer timer("part2")()
	result := 0
	// variables specific to this problem
	topBottomPattern := regexp.MustCompile(`^#+$`)
	gridRowPattern := regexp.MustCompile(`^#[.O@#]+#$`)
	instructionPattern := regexp.MustCompile(`[<v>^]+`)
	var grid [][]*Position
	var robotPosition Position
	fullInstructions := ""

	// Begin file parsing
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Day-specific code
		row := make([]*Position, 0)
		if topBottomPattern.MatchString(line) || gridRowPattern.MatchString(line) {
			for colIndex, char := range line {
				switch string(char) {
				case "#", ".":
					row = append(row, &Position{x: colIndex * 2, y: len(grid), content: char})
					row = append(row, &Position{x: colIndex*2 + 1, y: len(grid), content: char})
				case "O":
					row = append(row, &Position{x: colIndex * 2, y: len(grid), content: rune('[')})
					row = append(row, &Position{x: colIndex*2 + 1, y: len(grid), content: rune(']')})
				case "@":
					row = append(row, &Position{x: colIndex * 2, y: len(grid), content: char})
					robotPosition = Position{x: colIndex * 2, y: len(grid), content: char} // make a copy
					row = append(row, &Position{x: colIndex*2 + 1, y: len(grid), content: rune('.')})
				}
			}
			grid = append(grid, row)
		} else if instructionPattern.MatchString(line) {
			fullInstructions += line
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file: ", err)
	}
	// Post file-processing code.
	printGrid(grid)
	// Go through each instruction
	fmt.Println("Robot Position: ", robotPosition)
	for _, instruction := range fullInstructions {
		processInstruction2(grid, &robotPosition, instruction)
		// printGrid(grid)
	}
	result = calculateScore2(grid)
	fmt.Println("The final result is: ", result)
}
