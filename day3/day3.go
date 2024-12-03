package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func main() {
	part1()
	part2()
}

func calculate(line string) int {
	result := 0
	regexPattern := regexp.MustCompile(`mul\((\d+),(\d+)\)`)
	// find all the matches in the string
	matches := regexPattern.FindAllStringSubmatch(line, -1)
	for _, match := range matches {
		// convert the strings to integers
		first, err1 := strconv.Atoi(match[1])
		second, err2 := strconv.Atoi(match[2])
		if err1 != nil || err2 != nil {
			fmt.Println("Error converting string to int:", err1, err2)
			continue
		}
		// multiply the two numbers
		result += first * second
	}
	return result
}

func part1() {
	// https://adventofcode.com/2024/day/3
	// Regex and look for mul(3,4) style strings
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	result := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		result += calculate(line)
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file: ", err)
	}
	fmt.Println("The final result is: ", result)
}

func part2() {
	// https://adventofcode.com/2024/day/3#part2
	// https://adventofcode.com/2024/day/3
	// Regex and look for mul(3,4) style strings
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	doDontPattern := regexp.MustCompile(`(do|don't)\(\)`)
	scanner := bufio.NewScanner(file)
	oneline := ""
	for scanner.Scan() {
		line := scanner.Text()
		oneline += line
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file: ", err)
	}
	result := 0
	sections := doDontPattern.FindAllStringIndex(oneline, -1)
	for i := 0; i < len(sections); i++ {
		sectionIndex := sections[i]
		lastSectionIndex := []int{len(oneline), len(oneline)}
		nextSectionIndex := lastSectionIndex
		if (i + 1) < len(sections) {
			nextSectionIndex = sections[i+1]
		}
		if i == 0 {
			// First time, we need to add the first section (enabled by default)
			result += calculate(oneline[:sectionIndex[0]])
		}
		// convert the strings to integers
		section := oneline[sectionIndex[0]:sectionIndex[1]]
		if section == "do()" {
			result += calculate(oneline[(sectionIndex[0] + 4):nextSectionIndex[0]])
		}
	}
	fmt.Println("The final result is: ", result)
}
