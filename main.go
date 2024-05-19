package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func NewGraph() *Graph {
	return &Graph{rooms: make(map[string]*Room)}
}

func (g *Graph) AddRoom(name string, x, y int) {
	if _, exists := g.rooms[name]; !exists {
		g.rooms[name] = &Room{name: name, x: x, y: y}
	}
}

func (g *Graph) AddEdge(from, to string, capacity int) {
	fromRoom, fromExists := g.rooms[from]
	toRoom, toExists := g.rooms[to]
	if !fromExists || !toExists {
		fmt.Println("Error: Invalid data format, link to unknown room.")
		os.Exit(1)
	}
	for _, edge := range fromRoom.edges {
		if edge.to == toRoom {
			edge.capacity += capacity
			return
		}
	}
	fromEdge := &Edge{to: toRoom, capacity: capacity, flow: 0}
	toEdge := &Edge{to: fromRoom, capacity: 0, flow: 0}
	fromEdge.reverse = toEdge
	toEdge.reverse = fromEdge
	fromRoom.edges = append(fromRoom.edges, fromEdge)
	toRoom.edges = append(toRoom.edges, toEdge)
}

func input(filename string) (int, *Graph, string, string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error: File could not be opened.")
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	numAnts, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Println("Error: Ant count could not be obtained.")
		os.Exit(1)
	}
	graph := NewGraph()
	var startRoom, endRoom string
	lineNumber := 1
	for scanner.Scan() {
		line := scanner.Text()
		lineNumber++
		if len(line) == 0 || line[0] == '#' {
			if line == "##start" {
				if !scanner.Scan() {
					fmt.Println("Error: No start room found.")
					os.Exit(1)
				}
				lineNumber++
				startRoom = scanner.Text()
				parts := strings.Split(startRoom, " ")
				if len(parts) != 3 {
					fmt.Println("Error: Invalid start room format.")
					os.Exit(1)
				}
				x, _ := strconv.Atoi(parts[1])
				y, _ := strconv.Atoi(parts[2])
				graph.AddRoom(parts[0], x, y)
				startRoom = parts[0]
			} else if line == "##end" {
				if !scanner.Scan() {
					fmt.Println("Error: No end room found.")
					os.Exit(1)
				}
				lineNumber++
				endRoom = scanner.Text()
				parts := strings.Split(endRoom, " ")
				if len(parts) != 3 {
					fmt.Println("Error: Invalid end room format.")
					os.Exit(1)
				}
				x, _ := strconv.Atoi(parts[1])
				y, _ := strconv.Atoi(parts[2])
				graph.AddRoom(parts[0], x, y)
				endRoom = parts[0]
			}
			continue
		}
		if strings.Contains(line, "-") {
			parts := strings.Split(line, "-")
			if len(parts) != 2 {
				fmt.Println("Error: Invalid link format.")
				os.Exit(1)
			}
			graph.AddEdge(parts[0], parts[1], 1)
		} else {
			parts := strings.Split(line, " ")
			if len(parts) != 3 {
				fmt.Println("Error: Invalid room format.")
				os.Exit(1)
			}
			x, _ := strconv.Atoi(parts[1])
			y, _ := strconv.Atoi(parts[2])
			graph.AddRoom(parts[0], x, y)
		}
	}
	if startRoom == "" {
		fmt.Println("Error: No start room found.")
		os.Exit(1)
	}
	if endRoom == "" {
		fmt.Println("Error: No end room found.")
		os.Exit(1)
	}
	return numAnts, graph, startRoom, endRoom
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Terminal error. Usage: go run main.go <inputfile>")
		return
	}
	filename := os.Args[1]
	numAnts, graph, startRoom, endRoom := input(filename)
	maxFlow := edmondsKarp(graph, startRoom, endRoom)
	fmt.Printf("Max flow: %d\n", maxFlow)
	simAnts(numAnts, graph, startRoom, endRoom)
}
