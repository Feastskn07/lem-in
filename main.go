package main

import (
	"bufio"
	"container/list"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Room struct {
	name  string
	x, y  int
	edges []*Edge
}

type Edge struct {
	to       *Room
	capacity int
	flow     int
	reverse  *Edge
}

type Graph struct {
	rooms map[string]*Room
}

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

func edmondsKarp(graph *Graph, source, sink string) int {
	totalFlow := 0
	parent := make(map[*Room]*Edge)
	for BFS(graph, source, sink, parent) {
		pathFlow := int(^uint(0) >> 1)
		for v := graph.rooms[sink]; v.name != source; v = parent[v].reverse.to {
			pathFlow = min(pathFlow, parent[v].capacity-parent[v].flow)
		}
		for v := graph.rooms[sink]; v.name != source; v = parent[v].reverse.to {
			parent[v].flow += pathFlow
			parent[v].reverse.flow -= pathFlow
		}
		totalFlow += pathFlow
	}
	return totalFlow
}

func BFS(graph *Graph, source, sink string, parent map[*Room]*Edge) bool {
	visited := make(map[string]bool)
	queue := list.New()
	queue.PushBack(graph.rooms[source])
	visited[source] = true
	for queue.Len() > 0 {
		current := queue.Front().Value.(*Room)
		queue.Remove(queue.Front())
		for _, edge := range current.edges {
			if !visited[edge.to.name] && edge.flow < edge.capacity {
				visited[edge.to.name] = true
				parent[edge.to] = edge
				queue.PushBack(edge.to)
				if edge.to.name == sink {
					return true
				}
			}
		}
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func simAnts(numAnts int, graph *Graph, source, sink string) {
	type ant struct {
		id   int
		path []*Room
	}
	antPaths := make([][]*Room, numAnts)
	for i := 0; i < numAnts; i++ {
		antPaths[i] = []*Room{graph.rooms[source]}
	}
	steps := 0
	for {
		moveCount := 0
		steps++
		fmt.Printf("Step %d: \n", steps)
		for i := 0; i < numAnts; i++ {
			ant := antPaths[i]
			if len(ant) == 0 {
				continue
			}
			currentRoom := ant[len(ant)-1]
			if currentRoom.name == sink {
				continue
			}
			for _, edge := range currentRoom.edges {
				if edge.flow > 0 && edge.to.name != source {
					edge.flow--
					antPaths[i] = append(antPaths[i], edge.to)
					fmt.Printf("L%d-%s ", i+1, edge.to.name)
					moveCount++
					break
				}
			}
		}
		if moveCount == 0 {
			break
		}
		fmt.Println()
	}
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
