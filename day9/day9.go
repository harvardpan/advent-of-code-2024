package main

import (
	"bufio"
	"container/list"
	"fmt"
	"os"
	"strconv"
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

type Block struct {
	fileId int
	size   int
}

func findNextInsertIndex(startElement *list.Element, endElement *list.Element) *list.Element {
	for e := startElement; e != nil; e = e.Next() {
		if e.Value.(Block).fileId == -1 {
			return e
		}
		if e == endElement {
			break
		}
	}
	return nil
}

func calculateResult(disk *list.List) int {
	result := 0
	index := 0
	for e := disk.Front(); e != nil; e = e.Next() {
		if e.Value.(Block).fileId == -1 {
			continue
		}
		for j := 0; j < e.Value.(Block).size; j++ {
			result += e.Value.(Block).fileId * index
			index++
		}
	}
	return result
}

func printDisk(disk *list.List) {
	fmt.Print("Row: ")
	for e := disk.Front(); e != nil; e = e.Next() {
		if e.Value.(Block).fileId == -1 {
			fmt.Print(".")
		} else {
			fmt.Print(e.Value.(Block).fileId, "")
		}
	}
	fmt.Println()
}

func part1() {
	// https://adventofcode.com/2024/day/9
	// Walk an operations tree and determine a valid order of operations
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	defer timer("part1")()
	// result := 0
	// variables specific to this problem
	fileId := 0 // initial file ID that gets incremented
	disk := list.New()
	// Begin file parsing
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		for i, char := range line {
			size, _ := strconv.Atoi(string(char))
			if size == 0 {
				continue
			}
			if i%2 == 0 {
				// file definition
				disk.PushBack(Block{fileId, size})
				fileId++
			} else {
				// empty space
				disk.PushBack(Block{-1, size})
			}
		}
		printDisk(disk)
		break // only one row today
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file: ", err)
	}
	pInsertPosition := findNextInsertIndex(disk.Front(), disk.Back())
	if pInsertPosition == nil {
		fmt.Println("Completely filled.")
		return
	}
	for e := disk.Back(); e != nil; e = e.Prev() {
		if e == pInsertPosition {
			break
		}
		if e.Value.(Block).fileId == -1 {
			continue
		}
		// Need to move e to the insert position, but can only move to the size of the insert position
		// and then move the rest of the blocks to the next insert position
		if e.Value.(Block).size > pInsertPosition.Value.(Block).size {
			// Split up the "e" block, insert the first part at the insert position, and keep the rest for the next loop
			disk.InsertBefore(Block{e.Value.(Block).fileId, e.Value.(Block).size - pInsertPosition.Value.(Block).size}, e)
			pInsertPosition.Value = Block{e.Value.(Block).fileId, pInsertPosition.Value.(Block).size}
			e.Value = Block{-1, pInsertPosition.Value.(Block).size} // newly free space
		} else if e.Value.(Block).size < pInsertPosition.Value.(Block).size {
			// Split up the insert position block.
			disk.InsertAfter(Block{pInsertPosition.Value.(Block).fileId, pInsertPosition.Value.(Block).size - e.Value.(Block).size}, pInsertPosition)
			pInsertPosition.Value = Block{e.Value.(Block).fileId, e.Value.(Block).size}
			e.Value = Block{-1, e.Value.(Block).size}
		} else {
			pInsertPosition.Value = e.Value
			e.Value = Block{-1, e.Value.(Block).size}
		}
		pInsertPosition = findNextInsertIndex(pInsertPosition.Next(), disk.Back())
		if pInsertPosition == nil || pInsertPosition == e {
			break
		}
	}
	printDisk(disk)
	fmt.Println("The final result is: ", calculateResult(disk))
}

func part2() {
	// https://adventofcode.com/2024/day/9#part2
}
