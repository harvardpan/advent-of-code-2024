package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"math"
	"os"
	"time"
)

type Node struct {
	x       int
	y       int
	content rune
	prev    []*PriorityQueueItem // previous node in the path
}

// Priority Queue using the container/heap package.
// Took example code from https://pkg.go.dev/container/heap

// An PriorityQueueItem is something we manage in a priority queue.
type PriorityQueueItem struct {
	node      *Node
	direction int // north = 0, east = 90, south = 180, west = 270
	priority  int // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*PriorityQueueItem

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
	item := x.(*PriorityQueueItem)
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

func printGrid(grid [][]*Node, path []*PriorityQueueItem) {
	visited := make(map[string]string)
	for _, item := range path {
		directionString := ""
		switch item.direction {
		case 0:
			directionString = "^"
		case 90:
			directionString = ">"
		case 180:
			directionString = "v"
		case 270:
			directionString = "<"
		}
		visited[fmt.Sprintf("%d,%d", item.node.x, item.node.y)] = directionString
	}
	for rowIndex, row := range grid {
		for colIndex, item := range row {
			direction, exists := visited[fmt.Sprintf("%d,%d", colIndex, rowIndex)]
			if exists {
				fmt.Print(direction)
			} else {
				fmt.Print(string(item.content))
			}
		}
		fmt.Println()
	}
}

func constructShortestPath(item *PriorityQueueItem) []*PriorityQueueItem {
	path := make([]*PriorityQueueItem, 0)
	path = append(path, item)
	for item != nil && item.priority != 0 {
		minPriority := item.priority // the next step should always have a lower priority than the current item
		var minItem *PriorityQueueItem = nil
		for _, prevItem := range item.node.prev {
			if prevItem.priority < minPriority {
				minPriority = prevItem.priority
				minItem = prevItem
			}
		}
		if minItem == nil {
			break
		}
		path = append(path, minItem)
		item = minItem
	}
	return path
}

func constructShortestPaths(item *PriorityQueueItem) [][]*PriorityQueueItem {
	paths := make([][]*PriorityQueueItem, 0)
	if item == nil || item.priority == 0 {
		// terminating condition
		path := make([]*PriorityQueueItem, 0)
		path = append(path, item)
		paths = append(paths, path)
		return paths
	}

	minItems := make([]*PriorityQueueItem, 0)
	for _, prevItem := range item.node.prev {
		if prevItem.priority < item.priority {
			minItems = append(minItems, prevItem)
		}
	}
	for _, minItem := range minItems {
		newPaths := constructShortestPaths(minItem)
		for _, path := range newPaths {
			// prepend the current item to the path
			path = append([]*PriorityQueueItem{item}, path...)
			paths = append(paths, path)
		}
	}
	return paths
}

func findShortestPath(grid [][]*Node, pq PriorityQueue, visited map[string]bool, endPosition *Node, returnAllPaths bool) [][]*PriorityQueueItem {
	endItems := make([]*PriorityQueueItem, 0)
	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*PriorityQueueItem)
		if visited[fmt.Sprintf("%d,%d,%d", item.node.x, item.node.y, item.direction)] {
			continue
		}
		visited[fmt.Sprintf("%d,%d,%d", item.node.x, item.node.y, item.direction)] = true
		if item.node.x == endPosition.x && item.node.y == endPosition.y {
			fmt.Println("Found the end position")
			endItems = append(endItems, item)
			continue
		}
		// Add the next possible moves to the queue
		// Up
		if item.node.y-1 >= 0 && grid[item.node.y-1][item.node.x].content != rune('#') {
			additionalDistance := 0
			switch item.direction {
			case 0:
				additionalDistance = 1
			case 90, 270:
				additionalDistance = 1001 // 1000 for the turn, 1 for the move
			case 180:
				additionalDistance = 2001 // 2000 for 2 turns, 1 for the move
			}
			node := grid[item.node.y-1][item.node.x]
			node.prev = append(node.prev, item)
			heap.Push(&pq, &PriorityQueueItem{
				node:      node,
				direction: 0,
				priority:  item.priority + additionalDistance,
			})
		}
		// Down
		if item.node.y+1 < len(grid) && grid[item.node.y+1][item.node.x].content != rune('#') {
			additionalDistance := 0
			switch item.direction {
			case 180:
				additionalDistance = 1
			case 90, 270:
				additionalDistance = 1001 // 1000 for the turn, 1 for the move
			case 0:
				additionalDistance = 2001 // 2000 for 2 turns, 1 for the move
			}
			node := grid[item.node.y+1][item.node.x]
			node.prev = append(node.prev, item)
			heap.Push(&pq, &PriorityQueueItem{
				node:      node,
				direction: 180,
				priority:  item.priority + additionalDistance,
			})
		}
		// Left
		if item.node.x-1 >= 0 && grid[item.node.y][item.node.x-1].content != rune('#') {
			additionalDistance := 0
			switch item.direction {
			case 270:
				additionalDistance = 1
			case 0, 180:
				additionalDistance = 1001 // 1000 for the turn, 1 for the move
			case 90:
				additionalDistance = 2001 // 2000 for 2 turns, 1 for the move
			}
			node := grid[item.node.y][item.node.x-1]
			node.prev = append(node.prev, item)
			heap.Push(&pq, &PriorityQueueItem{
				node:      node,
				direction: 270,
				priority:  item.priority + additionalDistance,
			})
		}
		// Right
		if item.node.x+1 < len(grid[0]) && grid[item.node.y][item.node.x+1].content != rune('#') {
			additionalDistance := 0
			switch item.direction {
			case 90:
				additionalDistance = 1
			case 0, 180:
				additionalDistance = 1001 // 1000 for the turn, 1 for the move
			case 270:
				additionalDistance = 2001 // 2000 for 2 turns, 1 for the move
			}
			node := grid[item.node.y][item.node.x+1]
			node.prev = append(node.prev, item)
			heap.Push(&pq, &PriorityQueueItem{
				node:      node,
				direction: 90,
				priority:  item.priority + additionalDistance,
			})
		}
	}
	paths := make([][]*PriorityQueueItem, 0)
	minPriority := math.MaxInt64
	minItems := make([]*PriorityQueueItem, 0)
	for _, endItem := range endItems {
		if endItem.priority < minPriority {
			minPriority = endItem.priority
			minItems = make([]*PriorityQueueItem, 0)
			minItems = append(minItems, endItem)
		} else if endItem.priority == minPriority {
			minItems = append(minItems, endItem)
		}
	}
	for _, minItem := range minItems {
		if !returnAllPaths {
			return [][]*PriorityQueueItem{constructShortestPath(minItem)}
		}

		paths = append(paths, constructShortestPaths(minItem)...)
	}
	return paths
}

func part1() {
	// https://adventofcode.com/2024/day/16
	// Find shortest path through maze - Dijsktra's algorithm
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	defer timer("part1")()
	result := 0
	// variables specific to this problem
	grid := make([][]*Node, 0)
	var startPosition, endPosition *Node

	// Begin file parsing
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		row := make([]*Node, 0)
		// Day-specific code
		for colIndex, char := range line {
			node := &Node{x: colIndex, y: len(grid), content: char, prev: make([]*PriorityQueueItem, 0)}
			row = append(row, node)
			if char == 'S' {
				startPosition = node
			} else if char == 'E' {
				endPosition = node
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
	item := &PriorityQueueItem{
		node:      startPosition,
		direction: 90, // the start direction is east
		priority:  0,
	}
	heap.Push(&pq, item)
	visited := make(map[string]bool)
	path := findShortestPath(grid, pq, visited, endPosition, false)[0]
	if len(path) > 0 {
		result = path[0].priority
	}
	printGrid(grid, path)
	fmt.Println("The final result is: ", result)
}

func part2() {
	// https://adventofcode.com/2024/day/16#part2
	// Find all the shortest paths through the maze, and get the locations
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	defer timer("part2")()
	result := 0
	// variables specific to this problem
	grid := make([][]*Node, 0)
	var startPosition, endPosition *Node

	// Begin file parsing
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		row := make([]*Node, 0)
		// Day-specific code
		for colIndex, char := range line {
			node := &Node{x: colIndex, y: len(grid), content: char, prev: make([]*PriorityQueueItem, 0)}
			row = append(row, node)
			if char == 'S' {
				startPosition = node
			} else if char == 'E' {
				endPosition = node
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
	item := &PriorityQueueItem{
		node:      startPosition,
		direction: 90, // the start direction is east
		priority:  0,
	}
	heap.Push(&pq, item)
	visited := make(map[string]bool)
	paths := findShortestPath(grid, pq, visited, endPosition, true)
	uniqueLocations := make(map[string]bool)
	for _, path := range paths {
		for _, item := range path {
			uniqueLocations[fmt.Sprintf("%d,%d", item.node.x, item.node.y)] = true
		}
	}
	result = len(uniqueLocations)
	fmt.Println("The final result is: ", result)
}
