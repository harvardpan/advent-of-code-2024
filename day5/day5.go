package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
)

func main() {
	part1()
	part2()
}

func arePagesInOrder(pages []int, forwardMap map[int][]int, backwardMap map[int][]int) bool {
	// Loop through the pages and check the forwardMap and backwardMap
	for i, page := range pages {
		var pagesBehind []int
		var pagesAhead []int
		if i > 0 {
			pagesBehind = pages[:i]
		} else {
			pagesBehind = []int{}
		}
		if i+1 < len(pages) {
			pagesAhead = pages[i+1:]
		} else {
			pagesAhead = []int{}
		}
		for _, behind := range pagesBehind {
			// Look through the forwardMap. If it finds an entry there, then it is not in the right order
			if forwardMap[page] != nil {
				if slices.Contains(forwardMap[page], behind) {
					return false
				}
			}
		}
		for _, ahead := range pagesAhead {
			// Look through the backwardMap. If it finds an entry there, then it is not in the right order
			if backwardMap[page] != nil {
				if slices.Contains(backwardMap[page], ahead) {
					return false
				}
			}
		}
	}
	return true
}

func part1() {
	// https://adventofcode.com/2024/day/5
	// Confirm that pages are in the right order
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	forwardMap := make(map[int][]int)
	backwardMap := make(map[int][]int)
	orderPairPattern := regexp.MustCompile(`(\d+)\|(\d+)`)
	result := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		pages := make([]int, 0)
		if strings.Contains(line, "|") {
			parts := orderPairPattern.FindStringSubmatch(line)
			if len(parts) == 3 {
				left, _ := strconv.Atoi(parts[1])
				right, _ := strconv.Atoi(parts[2])
				forwardMap[left] = append(forwardMap[left], right)
				backwardMap[right] = append(backwardMap[right], left)
			}
		} else if strings.Contains(line, ",") {
			// Split the line by commas
			parts := strings.Split(line, ",")
			for _, part := range parts {
				pageNumber, _ := strconv.Atoi(part)
				pages = append(pages, pageNumber)
			}
			// Loop through the pages and check the forwardMap and backwardMap
			if arePagesInOrder(pages, forwardMap, backwardMap) {
				// Now we know this sequence of pages is in the right order.
				fmt.Println("Safety manual update (", line, ") is in the right order.")
				// Get the middle element
				result += pages[len(pages)/2]
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file: ", err)
	}
	fmt.Println("The final result is: ", result)
}

func part2() {
	// https://adventofcode.com/2024/day/5#part2
	// Confirm that pages are in the right order
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	forwardMap := make(map[int][]int)
	backwardMap := make(map[int][]int)
	orderPairPattern := regexp.MustCompile(`(\d+)\|(\d+)`)
	result := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		pages := make([]int, 0)
		if strings.Contains(line, "|") {
			parts := orderPairPattern.FindStringSubmatch(line)
			if len(parts) == 3 {
				left, _ := strconv.Atoi(parts[1])
				right, _ := strconv.Atoi(parts[2])
				forwardMap[left] = append(forwardMap[left], right)
				backwardMap[right] = append(backwardMap[right], left)
			}
		} else if strings.Contains(line, ",") {
			// Split the line by commas
			parts := strings.Split(line, ",")
			for _, part := range parts {
				pageNumber, _ := strconv.Atoi(part)
				pages = append(pages, pageNumber)
			}
			// Loop through the pages and check the forwardMap and backwardMap
			if !arePagesInOrder(pages, forwardMap, backwardMap) {
				// Now we know this sequence of pages is not in the right order.
				fmt.Println("Safety manual update (", line, ") is not in the right order.")
				// We have to sort the pages now based on a custom comparator function
				sort.Slice(pages, func(i, j int) bool {
					if forwardMap[pages[i]] != nil {
						if slices.Contains(forwardMap[pages[i]], pages[j]) {
							return true
						}
					}
					if backwardMap[pages[j]] != nil {
						if slices.Contains(backwardMap[pages[j]], pages[i]) {
							return true
						}
					}
					return false
				})
				fmt.Println("The sorted pages are: ", pages)
				// Get the middle element
				result += pages[len(pages)/2]
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file: ", err)
	}
	fmt.Println("The final result is: ", result)
}
