package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"time"
)

type Position struct {
	x, y      int
	direction int // degrees, with 0 being north, and 180 being south
}

// Priority Queue using the container/heap package.
// Took example code from https://pkg.go.dev/container/heap

// An Item is something we manage in a priority queue.
type Item struct {
	position     Position // The value of the item; arbitrary.
	previousItem *Item
	priority     int // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the lowest, not highest, priority so we use greater than here.
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // don't stop the GC from reclaiming the item eventually
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

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

func printGrid(grid [][]rune, path []*Item) {
	visited := make(map[string]string)
	for _, item := range path {
		directionString := ""
		switch item.position.direction {
		case 0:
			directionString = "^"
		case 90:
			directionString = ">"
		case 180:
			directionString = "v"
		case 270:
			directionString = "<"
		}
		visited[fmt.Sprintf("%d,%d", item.position.x, item.position.y)] = directionString
	}
	for rowIndex, row := range grid {
		for colIndex, char := range row {
			direction, exists := visited[fmt.Sprintf("%d,%d", colIndex, rowIndex)]
			if exists {
				fmt.Print(direction)
			} else {
				fmt.Print(string(char))
			}
		}
		fmt.Println()
	}
}

func findShortestPath(grid [][]rune, pq PriorityQueue, visited map[string]bool, endPosition Position) []*Item {
	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Item)
		if visited[fmt.Sprintf("%d,%d", item.position.x, item.position.y)] {
			continue
		}
		visited[fmt.Sprintf("%d,%d", item.position.x, item.position.y)] = true
		if item.position.x == endPosition.x && item.position.y == endPosition.y {
			fmt.Println("Found the end position")
			// Build up the path by checking previousItem for each item
			path := make([]*Item, 0)
			for item != nil {
				path = append(path, item)
				item = item.previousItem
			}
			return path
		}
		// Add the next possible moves to the queue
		// Up
		if item.position.y-1 >= 0 && grid[item.position.y-1][item.position.x] != '#' {
			additionalDistance := 0
			switch item.position.direction {
			case 0:
				additionalDistance = 1
			case 90, 270:
				additionalDistance = 1001 // 1000 for the turn, 1 for the move
			case 180:
				additionalDistance = 2001 // 2000 for 2 turns, 1 for the move
			}
			heap.Push(&pq, &Item{
				position:     Position{x: item.position.x, y: item.position.y - 1, direction: 0},
				priority:     item.priority + additionalDistance,
				previousItem: item,
			})
		}
		// Down
		if item.position.y+1 < len(grid) && grid[item.position.y+1][item.position.x] != '#' {
			additionalDistance := 0
			switch item.position.direction {
			case 180:
				additionalDistance = 1
			case 90, 270:
				additionalDistance = 1001 // 1000 for the turn, 1 for the move
			case 0:
				additionalDistance = 2001 // 2000 for 2 turns, 1 for the move
			}
			heap.Push(&pq, &Item{
				position:     Position{x: item.position.x, y: item.position.y + 1, direction: 180},
				priority:     item.priority + additionalDistance,
				previousItem: item,
			})
		}
		// Left
		if item.position.x-1 >= 0 && grid[item.position.y][item.position.x-1] != '#' {
			additionalDistance := 0
			switch item.position.direction {
			case 270:
				additionalDistance = 1
			case 0, 180:
				additionalDistance = 1001 // 1000 for the turn, 1 for the move
			case 90:
				additionalDistance = 2001 // 2000 for 2 turns, 1 for the move
			}
			heap.Push(&pq, &Item{
				position:     Position{x: item.position.x - 1, y: item.position.y, direction: 270},
				priority:     item.priority + additionalDistance,
				previousItem: item,
			})
		}
		// Right
		if item.position.x+1 < len(grid[0]) && grid[item.position.y][item.position.x+1] != '#' {
			additionalDistance := 0
			switch item.position.direction {
			case 90:
				additionalDistance = 1
			case 0, 180:
				additionalDistance = 1001 // 1000 for the turn, 1 for the move
			case 270:
				additionalDistance = 2001 // 2000 for 2 turns, 1 for the move
			}
			heap.Push(&pq, &Item{
				position:     Position{x: item.position.x + 1, y: item.position.y, direction: 90},
				priority:     item.priority + additionalDistance,
				previousItem: item,
			})
		}
	}
	return make([]*Item, 0)
}

func part1() {
	// https://adventofcode.com/2024/day/16
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
	grid := make([][]rune, 0)
	var startPosition, endPosition Position

	// Begin file parsing
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		row := make([]rune, 0)
		// Day-specific code
		for colIndex, char := range line {
			row = append(row, char)
			if char == 'S' {
				startPosition = Position{x: colIndex, y: len(grid), direction: 90} // start facing east
			} else if char == 'E' {
				endPosition = Position{x: colIndex, y: len(grid)}
			}
		}
		grid = append(grid, row)
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file: ", err)
	}
	// Post file-processing code.
	pq := make(PriorityQueue, 0) // priority queue makes dijkstra much easier
	heap.Init(&pq)
	// Add the start node
	item := &Item{
		position: startPosition,
		priority: 0,
	}
	heap.Push(&pq, item)
	visited := make(map[string]bool)
	path := findShortestPath(grid, pq, visited, endPosition)
	if len(path) > 0 {
		result = path[0].priority
	}
	printGrid(grid, path)
	fmt.Println("The final result is: ", result)
}

func part2() {
	// https://adventofcode.com/2024/day/16#part2
	//
}
