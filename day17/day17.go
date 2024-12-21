package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
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

// Global Variables
var registerA, registerB, registerC int

func getComboValue(operand string) int {
	// Combo operands 0 through 3 represent literal values 0 through 3.
	// Combo operand 4 represents the value of register A.
	// Combo operand 5 represents the value of register B.
	// Combo operand 6 represents the value of register C.
	// Combo operand 7 is reserved and will not appear in valid programs.
	switch operand {
	case "0", "1", "2", "3":
		value, _ := strconv.Atoi(operand)
		return value
	case "4":
		return registerA
	case "5":
		return registerB
	case "6":
		return registerC
	}
	return -1
}

func getLiteralValue(operand string) int {
	value, _ := strconv.Atoi(operand)
	return value
}

func processInstruction(instructions []string, i int) (int, string) {
	// Post file-processing code.
	opcode := instructions[i]
	switch opcode {
	case "0": // "adv" - division. numerator is register A, denominator is 2^(combo operand value)
		// The result of the division operation is truncated to an integer and then written to the A register.
		registerA = registerA / (1 << getComboValue(instructions[i+1]))
		return i + 2, ""
	case "1": // "bxl"
		// bitwise XOR of register B and the instruction's literal operand, then stores the result in register B.
		registerB = registerB ^ getLiteralValue(instructions[i+1])
		return i + 2, ""
	case "2": // "bst"
		// value of its combo operand modulo 8 (thereby keeping only its lowest 3 bits), then writes that value to the B register.
		registerB = getComboValue(instructions[i+1]) % 8
		return i + 2, ""
	case "3": // "jnz"
		// does nothing if the A register is 0. However, if the A register is not zero, it jumps by setting the instruction pointer to the value of its literal operand; if this instruction jumps, the instruction pointer is not increased by 2 after this instruction.
		if registerA != 0 {
			i = getLiteralValue(instructions[i+1])
			return i, ""
		} else {
			return i + 2, ""
		}
	case "4": // "bxc"
		// bitwise XOR of register B and register C, then stores the result in register B.
		registerB = registerB ^ registerC
		return i + 2, ""
	case "5": // "out"
		// calculates the value of its combo operand modulo 8, then outputs that value.
		output := strconv.Itoa(getComboValue(instructions[i+1]) % 8)
		return i + 2, output
	case "6": // "bdv"
		// same as adv, except result stored in register B
		registerB = registerA / (1 << getComboValue(instructions[i+1]))
		return i + 2, ""
	case "7": // "cdv"
		// same as adv, except result stored in register C
		registerC = registerA / (1 << getComboValue(instructions[i+1]))
		return i + 2, ""
	}
	return -1, ""
}

func part1() {
	// https://adventofcode.com/2024/day/17
	//
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	defer timer("part1")()
	result := 0
	// variables specific to this problem
	registerPattern := regexp.MustCompile(`^Register (\w): (\d+)$`)
	instructionPattern := regexp.MustCompile(`^Program: ([\d,]+)$`)
	var instructions []string
	outputs := make([]string, 0)

	// Begin file parsing
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Day-specific code
		if registerPattern.MatchString(line) {
			matches := registerPattern.FindStringSubmatch(line)
			register := matches[1]
			value, _ := strconv.Atoi(matches[2])
			switch register {
			case "A":
				registerA = value
			case "B":
				registerB = value
			case "C":
				registerC = value
			}
		} else if instructionPattern.MatchString(line) {
			matches := instructionPattern.FindStringSubmatch(line)
			instructions = strings.Split(matches[1], ",")
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file: ", err)
	}
	// Post file-processing code.
	for i := 0; i < len(instructions); {
		var output string
		i, output = processInstruction(instructions, i)
		if output != "" {
			outputs = append(outputs, output)
		}
	}
	// Joins the outputs into a single string, separated by commas
	fmt.Println("Output: ", strings.Join(outputs, ","))
	fmt.Println("The final result is: ", result)
}

func part2() {
	// https://adventofcode.com/2024/day/17#part2
	// calculate the value of register A that will output the program that was input originally
}
