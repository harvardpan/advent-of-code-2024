package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	part1()
	part2()
}

func part1() {
	// https://adventofcode.com/2024/day/1
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// declare the arrays that will store the values
	var firstnumbers []int
	var secondnumbers []int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// split the text by whitespace
		tokens := strings.Fields(line)
		if len(tokens) < 2 {
			continue
		}
		firstnum, err1 := strconv.Atoi(tokens[0])
		secondnum, err2 := strconv.Atoi(tokens[1])
		if err1 != nil || err2 != nil {
			fmt.Println("Error converting string to int")
			continue
		}
		firstnumbers = append(firstnumbers, firstnum)
		secondnumbers = append(secondnumbers, secondnum)
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	distance := 0
	// sort the arrays
	sort.Ints(firstnumbers)
	sort.Ints(secondnumbers)
	for i := 0; i < len(firstnumbers); i++ {
		firstnum := firstnumbers[i]
		secondnum := secondnumbers[i]
		// take the absolute value of the difference
		distance += int(math.Abs(float64(firstnum - secondnum)))
	}
	fmt.Println("The distance is:", distance)
}

func part2() {
	// https://adventofcode.com/2024/day/1#part2
	// Calculate a similarity score. Multiply the number on left with the number
	// of times that it appears on the right. Add up all the scores.
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// declare the arrays that will store the values
	firstnumbers := make([]int, 0)
	counts := make(map[int]int) // count the number of times the number appears in the second list
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// split the text by whitespace
		tokens := strings.Fields(line)
		if len(tokens) < 2 {
			continue
		}
		firstnum, err1 := strconv.Atoi(tokens[0])
		secondnum, err2 := strconv.Atoi(tokens[1])
		if err1 != nil || err2 != nil {
			fmt.Println("Error converting string to int")
			continue
		}
		firstnumbers = append(firstnumbers, firstnum)
		counts[secondnum] = counts[secondnum] + 1
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	similarity := 0
	for i := 0; i < len(firstnumbers); i++ {
		firstnum := firstnumbers[i]
		similarity += firstnum * counts[firstnum]
	}
	fmt.Println("The similarity is:", similarity)
}
