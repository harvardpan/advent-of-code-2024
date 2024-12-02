package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	part1()
	part2()
}

func isSafeReport(levels []int) bool {
	// assumes there are no negative numbers for levels
	previousLevel := -1
	direction := 0
	for _, level := range levels {
		if previousLevel < 0 {
			previousLevel = level
			continue
		}
		if level == previousLevel {
			// same level, so cannot be increasing or decreasing
			return false
		}
		var currentDirection int
		if level > previousLevel {
			currentDirection = 1
		} else if level < previousLevel {
			currentDirection = -1
		}
		if direction != 0 && currentDirection != direction {
			return false
		}
		// Direction check is good. Now check to see if they differ between 1-3 levels.
		if math.Abs(float64(level-previousLevel)) > 3 {
			return false
		}
		direction = currentDirection
		previousLevel = level
	}
	return true
}

func part1() {
	// https://adventofcode.com/2024/day/2
	// Calculate the number of "safe" levels. A level is safe if the levels are either increasing or decreasing by 1-3 levels.
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	numSafeReports := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// split the text by whitespace
		tokens := strings.Fields(line)
		// Convert the tokens to integers to get the list of levels
		levels := make([]int, 0)
		for _, token := range tokens {
			level, err := strconv.Atoi(token)
			if err != nil {
				fmt.Println("Error converting string to int:", err)
				continue
			}
			levels = append(levels, level)
		}
		if isSafeReport(levels) {
			numSafeReports++
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file: ", err)
	}
	fmt.Println("The # of safe reports is: ", numSafeReports)
}

func safeAfterDampening(levels []int) bool {
	for i := 0; i < len(levels); i++ {
		dampenedLevels := make([]int, 0)
		if i > 0 {
			dampenedLevels = append(dampenedLevels, levels[:i]...)
		}
		if i < len(levels)-1 {
			dampenedLevels = append(dampenedLevels, levels[i+1:]...)
		}
		if isSafeReport(dampenedLevels) {
			return true
		}
	}
	return false
}

func part2() {
	// https://adventofcode.com/2024/day/2
	// Calculate the number of "safe" levels. A level is safe if the levels are either increasing or decreasing by 1-3 levels.
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	numSafeReports := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// split the text by whitespace
		tokens := strings.Fields(line)
		// Convert the tokens to integers to get the list of levels
		levels := make([]int, 0)
		for _, token := range tokens {
			level, err := strconv.Atoi(token)
			if err != nil {
				fmt.Println("Error converting string to int:", err)
				continue
			}
			levels = append(levels, level)
		}
		if isSafeReport(levels) {
			numSafeReports++
		} else {
			if safeAfterDampening(levels) {
				numSafeReports++
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file: ", err)
	}
	fmt.Println("The # of safe reports is: ", numSafeReports)
}
