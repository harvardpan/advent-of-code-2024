package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", name, time.Since(start))
	}
}

type Plot struct {
	numWalls      int
	plant         rune
	x             int
	y             int
	neighbors     []*Plot
	soloCorners   int
	sharedCorners int
}

func main() {
	part1()
	part2()
}

func safeGetNeighbor(grid [][]*Plot, plot *Plot, neighbor [2]int) *Plot {
	if plot.x+neighbor[0] < 0 || plot.x+neighbor[0] >= len(grid[0]) {
		return nil
	}
	if plot.y+neighbor[1] < 0 || plot.y+neighbor[1] >= len(grid) {
		return nil
	}
	adjacentPlot := grid[plot.y+neighbor[1]][plot.x+neighbor[0]]
	if adjacentPlot.plant != plot.plant {
		return nil
	}
	return adjacentPlot
}

func checkCorner(grid [][]*Plot, plot *Plot, dir1 [2]int, dir2 [2]int, dir3 [2]int) {
	if safeGetNeighbor(grid, plot, dir1) == nil && safeGetNeighbor(grid, plot, dir2) == nil {
		plot.soloCorners++
	} else if safeGetNeighbor(grid, plot, dir1) != nil && safeGetNeighbor(grid, plot, dir2) != nil && safeGetNeighbor(grid, plot, dir3) == nil {
		plot.sharedCorners++
	} else if safeGetNeighbor(grid, plot, dir1) != nil && safeGetNeighbor(grid, plot, dir3) != nil && safeGetNeighbor(grid, plot, dir2) == nil {
		plot.sharedCorners++
	} else if safeGetNeighbor(grid, plot, dir2) != nil && safeGetNeighbor(grid, plot, dir3) != nil && safeGetNeighbor(grid, plot, dir1) == nil {
		plot.sharedCorners++
	}
}

func updatePlot(grid [][]*Plot, plot *Plot) {
	// update the number of walls around the plot
	if plot.numWalls != -1 {
		return
	}
	n := [2]int{0, -1}
	e := [2]int{1, 0}
	s := [2]int{0, 1}
	w := [2]int{-1, 0}
	ne := [2]int{1, -1}
	se := [2]int{1, 1}
	sw := [2]int{-1, 1}
	nw := [2]int{-1, -1}
	// check the n, e, s, w surrounding plots
	neigbors := [4][2]int{n, e, s, w}
	plot.numWalls = 4 // start with 4 walls
	for _, neighbor := range neigbors {
		adjacentPlot := safeGetNeighbor(grid, plot, neighbor)
		if adjacentPlot == nil {
			continue
		}
		plot.numWalls--
		plot.neighbors = append(plot.neighbors, adjacentPlot)
	}

	// for part 2, calculate the number of corners.
	plot.soloCorners = 0
	plot.sharedCorners = 0
	// 1. If you have no neighbors, you have a corner.
	// 2. If you have two neighbors out of three in that corner, you have a corner.
	// Check ne corner.
	checkCorner(grid, plot, n, e, ne)
	// Check se corner.
	checkCorner(grid, plot, e, s, se)
	// Check sw corner.
	checkCorner(grid, plot, s, w, sw)
	// Check nw corner.
	checkCorner(grid, plot, w, n, nw)
}

func explorePlot(plot *Plot, visited map[*Plot]bool, area *int, perimeter *int) (int, int) {
	if visited[plot] {
		return 0, 0
	}
	visited[plot] = true
	(*area)++
	(*perimeter) += plot.numWalls
	for _, neighbor := range plot.neighbors {
		explorePlot(neighbor, visited, area, perimeter)
	}
	return *area, *perimeter
}

func part1() {
	// https://adventofcode.com/2024/day/12
	// Calculate plot areas and perimeters
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	defer timer("part1")()
	result := 0
	// variables specific to this problem
	var grid [][]*Plot
	plots := make(map[rune][]*Plot)
	rowIndex := 0
	// Begin file parsing
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Day-specific code
		var row []*Plot
		for colIndex, char := range line {
			plot := Plot{numWalls: -1, plant: char, x: colIndex, y: rowIndex, neighbors: make([]*Plot, 0)}
			row = append(row, &plot)
			plotList, exists := plots[char]
			if !exists {
				plotList = make([]*Plot, 0)
			}
			plotList = append(plotList, &plot)
			plots[char] = plotList
		}
		grid = append(grid, row)
		rowIndex++
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file: ", err)
	}
	// Post file-processing code.
	// Update the data in the Plot structs
	for _, plotList := range plots {
		for _, plot := range plotList {
			updatePlot(grid, plot)
		}
	}
	// Calculate the area and perimeter of each plot
	for _, plotList := range plots {
		visited := make(map[*Plot]bool)
		area := 0
		perimeter := 0
		for _, plot := range plotList {
			area, perimeter = explorePlot(plot, visited, &area, &perimeter)
			if area > 0 {
				// fmt.Printf("Plant: %c, Area: %d, Perimeter: %d\n", plant, area, perimeter)
				result += area * perimeter
				area = 0
				perimeter = 0
			}
		}
	}
	fmt.Println("The final result is: ", result)
}

func explorePlot2(plot *Plot, visited map[*Plot]bool, area *int, soloCorners *int, sharedCorners *int) (int, int, int) {
	if visited[plot] {
		return 0, 0, 0
	}
	visited[plot] = true
	(*area)++
	(*soloCorners) += plot.soloCorners
	(*sharedCorners) += plot.sharedCorners
	for _, neighbor := range plot.neighbors {
		explorePlot2(neighbor, visited, area, soloCorners, sharedCorners)
	}
	return *area, *soloCorners, *sharedCorners
}

func part2() {
	// https://adventofcode.com/2024/day/12#part2
	// Calculate number of sides instead of perimeter (i.e. detect straight lines)
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	defer timer("part2")()
	result := 0
	// variables specific to this problem
	var grid [][]*Plot
	plots := make(map[rune][]*Plot)
	rowIndex := 0
	// Begin file parsing
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Day-specific code
		var row []*Plot
		for colIndex, char := range line {
			plot := Plot{numWalls: -1, plant: char, x: colIndex, y: rowIndex, neighbors: make([]*Plot, 0)}
			row = append(row, &plot)
			plotList, exists := plots[char]
			if !exists {
				plotList = make([]*Plot, 0)
			}
			plotList = append(plotList, &plot)
			plots[char] = plotList
		}
		grid = append(grid, row)
		rowIndex++
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file: ", err)
	}
	// Post file-processing code.
	// Update the data in the Plot structs
	for _, plotList := range plots {
		for _, plot := range plotList {
			updatePlot(grid, plot)
		}
	}
	// Calculate the area and perimeter of each plot
	for _, plotList := range plots {
		visited := make(map[*Plot]bool)
		area := 0
		soloCorners := 0
		sharedCorners := 0
		for _, plot := range plotList {
			area, soloCorners, sharedCorners = explorePlot2(plot, visited, &area, &soloCorners, &sharedCorners)
			sides := soloCorners + sharedCorners/3
			if area > 0 {
				// fmt.Printf("Plant: %c, Area: %d, Sides: %d, Solo: %d, Shared: %d\n", plant, area, sides, soloCorners, sharedCorners)
				result += area * sides
				area = 0
				soloCorners = 0
				sharedCorners = 0
			}
		}
	}
	fmt.Println("The final result is: ", result)
}
