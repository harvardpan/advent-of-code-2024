package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"time"
)

func timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", name, time.Since(start))
	}
}

type Robot struct {
	posX, posY int
	velX, velY int
}

func main() {
	part1()
	part2()
}

func step(robots []*Robot, gridSizeX, gridSizeY int) {
	for i := range robots {
		// modulo, handling negative numbers
		robots[i].posX = ((robots[i].posX+robots[i].velX)%gridSizeX + gridSizeX) % gridSizeX
		robots[i].posY = ((robots[i].posY+robots[i].velY)%gridSizeY + gridSizeY) % gridSizeY
	}
}

func printGrid(robots []*Robot, gridSizeX, gridSizeY int) {
	grid := make([][]rune, gridSizeY)
	for i := range grid {
		grid[i] = make([]rune, gridSizeX)
		for j := range grid[i] {
			grid[i][j] = '.'
		}
	}
	for _, robot := range robots {
		grid[robot.posY][robot.posX] = '#'
	}
	for i := range grid {
		fmt.Println(string(grid[i]))
	}
}

func calculateSafetyFactor(robots []*Robot, gridSizeX, gridSizeY int) int {
	// Calculate the safety factor of the robots
	nw, ne, sw, se := 0, 0, 0, 0
	for _, robot := range robots {
		if robot.posX < gridSizeX/2 && robot.posY < gridSizeY/2 {
			nw++
		} else if robot.posX > gridSizeX/2 && robot.posY < gridSizeY/2 {
			ne++
		} else if robot.posX < gridSizeX/2 && robot.posY > gridSizeY/2 {
			sw++
		} else if robot.posX > gridSizeX/2 && robot.posY > gridSizeY/2 {
			se++
		}
		// ignore anything that's in the center cross
	}
	// fmt.Println("NW: ", nw, " NE: ", ne, " SW: ", sw, " SE: ", se)
	return nw * ne * sw * se
}

func part1() {
	// https://adventofcode.com/2024/day/14
	// Calculate positions of all robots after 100 steps and calculate the safety factor
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	defer timer("part1")()
	result := 0
	// variables specific to this problem
	pattern := regexp.MustCompile(`p=(-?\d+),(-?\d+) v=(-?\d+),(-?\d+)`)
	gridSizeX := 101
	gridSizeY := 103
	steps := 100
	robots := make([]*Robot, 0)

	// Begin file parsing
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Day-specific code
		if pattern.MatchString(line) {
			matches := pattern.FindStringSubmatch(line)
			// fmt.Println(matches)
			posX, err := strconv.Atoi(matches[1])
			if err != nil {
				fmt.Println("Error converting posX portion:", err)
				return
			}
			posY, err := strconv.Atoi(matches[2])
			if err != nil {
				fmt.Println("Error converting posY portion:", err)
				return
			}
			velX, err := strconv.Atoi(matches[3])
			if err != nil {
				fmt.Println("Error converting velX portion:", err)
				return
			}
			velY, err := strconv.Atoi(matches[4])
			if err != nil {
				fmt.Println("Error converting velY portion:", err)
				return
			}
			robots = append(robots, &Robot{posX, posY, velX, velY})
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file: ", err)
	}
	// Post file-processing code.
	for i := 0; i < steps; i++ {
		step(robots, gridSizeX, gridSizeY)
	}
	result = calculateSafetyFactor(robots, gridSizeX, gridSizeY)
	fmt.Println("The final result is: ", result)
}

func deepCopyGrid(robots []*Robot) []*Robot {
	deepCopy := make([]*Robot, len(robots))
	for i, robot := range robots {
		deepCopy[i] = &Robot{posX: robot.posX, posY: robot.posY, velX: robot.velX, velY: robot.velY}
	}
	return deepCopy
}

func part2() {
	// https://adventofcode.com/2024/day/14#part2
	// Find the Christmas Tree easter egg. Do it by minimizing the safety factor.
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	defer timer("part2")()
	result := 0
	// variables specific to this problem
	pattern := regexp.MustCompile(`p=(-?\d+),(-?\d+) v=(-?\d+),(-?\d+)`)
	gridSizeX := 101
	gridSizeY := 103
	steps := 100000
	robots := make([]*Robot, 0)

	// Begin file parsing
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Day-specific code
		if pattern.MatchString(line) {
			matches := pattern.FindStringSubmatch(line)
			// fmt.Println(matches)
			posX, err := strconv.Atoi(matches[1])
			if err != nil {
				fmt.Println("Error converting posX portion:", err)
				return
			}
			posY, err := strconv.Atoi(matches[2])
			if err != nil {
				fmt.Println("Error converting posY portion:", err)
				return
			}
			velX, err := strconv.Atoi(matches[3])
			if err != nil {
				fmt.Println("Error converting velX portion:", err)
				return
			}
			velY, err := strconv.Atoi(matches[4])
			if err != nil {
				fmt.Println("Error converting velY portion:", err)
				return
			}
			robots = append(robots, &Robot{posX, posY, velX, velY})
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file: ", err)
	}
	// Post file-processing code.
	minSafetyFactor := math.MaxInt64
	var minSafetyGrid []*Robot
	minSafetyFactorStep := 0
	for i := 0; i < steps; i++ {
		step(robots, gridSizeX, gridSizeY)
		safetyFactor := calculateSafetyFactor(robots, gridSizeX, gridSizeY)
		if safetyFactor < minSafetyFactor {
			minSafetyFactor = safetyFactor
			minSafetyFactorStep = i + 1
			minSafetyGrid = deepCopyGrid(robots)
		}
	}
	printGrid(minSafetyGrid, gridSizeX, gridSizeY)
	fmt.Println("The step with the minimum safety factor is: ", minSafetyFactorStep)
	fmt.Println("The final result is: ", result)
}
