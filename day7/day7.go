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

type Operation struct {
	operand      int
	operator     string
	currentValue int
}

type Stack struct {
	items []Operation
}

func (stack Stack) IsEmpty() bool {
	return len(stack.items) == 0
}

func (s *Stack) Push(data Operation) {
	s.items = append(s.items, data)
}

func (s *Stack) Pop() Operation {
	if s.IsEmpty() {
		return Operation{}
	}
	popped := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return popped
}

func printSolution(stack Stack, firstValue int) {
	// Print the solution
	fmt.Print(firstValue, " ")
	for i := len(stack.items) - 1; i >= 0; i-- {
		fmt.Print(stack.items[i].operator, " ", stack.items[i].operand, " ")
	}
	fmt.Println("=", stack.items[0].currentValue)
}

func calculate(testValue int, operands []int) bool {
	// Traverse the operands backwards and see if we can get to the testValue
	// If we can, return true. Otherwise, return false
	stack := Stack{items: make([]Operation, 0)}
	currentValue := testValue // initial value is the testValue
	nextOperator := "*"
	index := len(operands) - 1
	for {
		operand := operands[index]
		// two possible operators: + and *, so we need to use the opposite operator - and /
		// we start with the division and see if we can get a valid whole number
		newValue := currentValue
		if currentValue%operand == 0 && nextOperator == "*" {
			newValue = currentValue / operand
			stack.Push(Operation{operand: operand, operator: "*", currentValue: currentValue})
		} else {
			newValue = currentValue - operand
			stack.Push(Operation{operand: operand, operator: "+", currentValue: currentValue})
			nextOperator = "*" // reset so it'll go down the first path next time
		}
		index--
		if index == 0 {
			// We are at the beginning.
			if newValue == operands[0] {
				printSolution(stack, newValue)
				return true
			}
			for {
				// Pop the stack until we find a "*" operator IFF there is at least one "*" operator
				popped := stack.Pop()
				if popped == (Operation{}) {
					// Popped all the way to the top, so we are done.
					return false
				}
				newValue = popped.currentValue
				index++
				if popped.operator == "*" {
					// We switch the operator to "+" and try the next path
					nextOperator = "+"
					break
				}
			}
		}
		currentValue = newValue
	}
}

func part1() {
	// https://adventofcode.com/2024/day/7
	// Walk an operations tree and determine a valid order of operations
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	defer timer("part1")()
	result := 0
	linePattern := regexp.MustCompile(`(\d+):(( \d+)+)`)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		matches := linePattern.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			testValue, _ := strconv.Atoi(match[1])
			// split and trim the operands
			operandsString := strings.Split(strings.TrimSpace(match[2]), " ")
			operands := make([]int, len(operandsString))
			for i, operand := range operandsString {
				operands[i], _ = strconv.Atoi(strings.TrimSpace(operand))
			}
			if calculate(testValue, operands) {
				result += testValue
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file: ", err)
	}
	// We should have the grid and the starting position now. Let's navigate the grid
	fmt.Println("The final result is: ", result)
}

func popUntilNextPath(stack *Stack, currentIndex int) (newValue int, nextOperator string, newIndex int, finished bool) {
	newIndex = currentIndex
	for {
		// Pop the stack until we find only "||" operators left
		popped := stack.Pop()
		if popped == (Operation{}) {
			// Popped all the way to the top, so we are done.
			finished = true
			break
		}
		newValue = popped.currentValue
		newIndex++
		if popped.operator == "*" {
			// We switch the operator to "+" and try the next path
			nextOperator = "+"
			break
		} else if popped.operator == "+" {
			// We switch the operator to "||" and try the next path
			nextOperator = "||"
			break
		}
	}
	return newValue, nextOperator, newIndex, finished
}

func calculate2(testValue int, operands []int) bool {
	// Traverse the operands backwards and see if we can get to the testValue
	// If we can, return true. Otherwise, return false
	stack := Stack{items: make([]Operation, 0)}
	currentValue := testValue // initial value is the testValue
	nextOperator := "*"
	index := len(operands) - 1
	for {
		operand := operands[index]
		// three possible operators: + and * and ||, so we need to use the opposite operator - and / and ||
		completedOperation := false
		newValue := currentValue
		if nextOperator == "*" {
			nextOperator = "+"
			if currentValue%operand == 0 {
				newValue = currentValue / operand
				stack.Push(Operation{operand: operand, operator: "*", currentValue: currentValue})
				completedOperation = true
			}
		}
		if !completedOperation && nextOperator == "+" {
			newValue = currentValue - operand
			stack.Push(Operation{operand: operand, operator: "+", currentValue: currentValue})
			nextOperator = "||"
			completedOperation = true
		}

		if !completedOperation && nextOperator == "||" {
			if strings.HasSuffix(strconv.Itoa(currentValue), strconv.Itoa(operand)) {
				newValue, _ = strconv.Atoi(strings.TrimSuffix(strconv.Itoa(currentValue), strconv.Itoa(operand)))
				stack.Push(Operation{operand: operand, operator: "||", currentValue: currentValue})
			} else {
				var finished bool
				newValue, nextOperator, index, finished = popUntilNextPath(&stack, index)
				if finished {
					return false
				}
				currentValue = newValue
				continue
			}
		}

		index--
		nextOperator = "*"
		if index == 0 {
			// We are at the beginning.
			if newValue == operands[0] {
				printSolution(stack, newValue)
				return true
			}
			var finished bool
			newValue, nextOperator, index, finished = popUntilNextPath(&stack, index)
			if finished {
				return false
			}
		}
		currentValue = newValue
	}
}

func part2() {
	// https://adventofcode.com/2024/day/7#part2
	// Walk an operations tree and determine a valid order of operations
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	defer timer("part2")()
	result := 0
	linePattern := regexp.MustCompile(`(\d+):(( \d+)+)`)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		matches := linePattern.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			testValue, _ := strconv.Atoi(match[1])
			// split and trim the operands
			operandsString := strings.Split(strings.TrimSpace(match[2]), " ")
			operands := make([]int, len(operandsString))
			for i, operand := range operandsString {
				operands[i], _ = strconv.Atoi(strings.TrimSpace(operand))
			}
			if calculate2(testValue, operands) {
				result += testValue
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file: ", err)
	}
	// We should have the grid and the starting position now. Let's navigate the grid
	fmt.Println("The final result is: ", result)
}
