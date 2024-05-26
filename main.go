package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Room struct {
	name        string
	visited     bool
	connections []string
	antsInside  int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Lütfen bir dosya ismi girin: go run main.go <dosya_ismi>")
		return
	}

	filename := os.Args[1]

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Dosya açılamadı:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var totalAnts int
	rooms := make(map[string]*Room)
	var startRoom, endRoom string
	readingRooms := false
	isStartRoom := false
	isEndRoom := false
	corridors := false

	lineCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		lineCount++

		if len(parts) == 0 {
			continue
		}

		if totalAnts == 0 {
			totalAnts, err = strconv.Atoi(parts[0])
			if err != nil {
				fmt.Println("Karınca sayısı okunamadı:", err)
				return
			}
			continue
		}

		if parts[0] == "##start" {
			isStartRoom = true
			readingRooms = true
			continue
		} else if parts[0] == "##end" {
			isEndRoom = true
			readingRooms = true
			continue
		}

		if readingRooms {
			if len(parts) != 3 {
				fmt.Printf("Hatalı oda formatı: %s, satır: %d\n", line, lineCount)
				return
			}
			roomName := parts[0]
			rooms[roomName] = &Room{name: roomName, visited: false, connections: []string{}, antsInside: 0}
			if isStartRoom {
				startRoom = roomName
				isStartRoom = false
			} else if isEndRoom {
				endRoom = roomName
				isEndRoom = false
				corridors = true
				readingRooms = false
			}
			continue
		}

		if corridors {
			if len(parts) == 1 && strings.Contains(parts[0], "-") {
				corridorParts := strings.Split(parts[0], "-")
				if len(corridorParts) != 2 {
					fmt.Printf("Hatalı koridor formatı: %s, satır: %d\n", line, lineCount)
					return
				}
				fromRoom := corridorParts[0]
				toRoom := corridorParts[1]
				if _, ok := rooms[fromRoom]; !ok {
					fmt.Printf("Bilinmeyen oda: %s, satır: %d\n", fromRoom, lineCount)
					return
				}
				if _, ok := rooms[toRoom]; !ok {
					fmt.Printf("Bilinmeyen oda: %s, satır: %d\n", toRoom, lineCount)
					return
				}
				rooms[fromRoom].connections = append(rooms[fromRoom].connections, toRoom)
				rooms[toRoom].connections = append(rooms[toRoom].connections, fromRoom)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Hata:", err)
	}

	// Find shortest path
	shortestPath := findShortestPath(rooms, startRoom, endRoom)

	// Assign ants to paths
	antPaths := assignAntsToPaths(totalAnts, shortestPath, rooms)

	// Print results
	printResults(totalAnts, antPaths)
}

func findShortestPath(rooms map[string]*Room, start string, end string) []string {
	var path []string

	queue := [][]string{{start}}

	for len(queue) > 0 {
		path = queue[0]
		queue = queue[1:]

		lastRoom := path[len(path)-1]

		if lastRoom == end {
			break
		}

		for _, connection := range rooms[lastRoom].connections {
			if !contains(path, connection) {
				newPath := append([]string{}, path...)
				newPath = append(newPath, connection)
				queue = append(queue, newPath)
			}
		}
	}

	return path
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func assignAntsToPaths(totalAnts int, shortestPath []string, rooms map[string]*Room) map[int][]string {
	antPaths := make(map[int][]string)

	// Initialize ants at the start room
	for ant := 1; ant <= totalAnts; ant++ {
		antPaths[ant] = []string{shortestPath[0]}
	}

	// Calculate next room for each ant based on available space in rooms
	for i := 1; i < len(shortestPath)-1; i++ {
		currentRoom := shortestPath[i]
		nextRoom := shortestPath[i+1]

		// Calculate total ants and rooms in the current room
		totalAntsInCurrentRoom := 0
		totalRoomsInCurrentRoom := len(rooms[currentRoom].connections) + 1 // including itself
		for _, conn := range rooms[currentRoom].connections {
			totalAntsInCurrentRoom += rooms[conn].antsInside
		}

		// Check if next room has enough space for the next ant
		if totalAntsInCurrentRoom < totalRoomsInCurrentRoom {
			for ant := i + 1; ant <= totalAnts; ant++ {
				antPaths[ant] = append(antPaths[ant], nextRoom)
				rooms[nextRoom].antsInside++
			}
		} else {
			// If next room doesn't have enough space, place the ants behind
			for ant := i + 1; ant <= totalAnts; ant++ {
				antPaths[ant] = append([]string{currentRoom}, antPaths[ant]...)
				rooms[currentRoom].antsInside++
			}
		}
	}

	return antPaths
}

func printResults(totalAnts int, antPaths map[int][]string) {
	// Print input file contents
	inputFile, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("Dosya açılamadı:", err)
		return
	}
	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Dosya okunamadı:", err)
		return
	}

	// Print a blank line
	fmt.Println()

	// Print ant paths
	for ant := 1; ant <= totalAnts; ant++ {
		fmt.Printf("L%d-%s\n", ant, strings.Join(antPaths[ant], " "))
	}
}
