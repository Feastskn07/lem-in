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
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Lütfen bir dosya ismi girin: go run main.go <dosya_ismi>")
		return
	}

	filename := os.Args[1]

	// Dosyayı aç
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
	var readingRooms bool = false

	lineCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		lineCount++

		// Boş satırları atla
		if len(parts) == 0 {
			continue
		}

		// İlk satırda karınca sayısını oku
		if totalAnts == 0 {
			totalAnts, err = strconv.Atoi(parts[0])
			if err != nil {
				fmt.Println("Karınca sayısı okunamadı:", err)
				return
			}
			continue
		}

		// ##start ve ##end işaretlerini kontrol et
		if parts[0] == "##start" {
			readingRooms = true
			continue
		} else if parts[0] == "##end" {
			readingRooms = true
			continue
		}

		// Odaları oku
		if readingRooms {
			if len(parts) != 3 {
				fmt.Printf("Hatalı oda formatı: %s, satır: %d\n", line, lineCount)
				return
			}
			roomName := parts[0]
			rooms[roomName] = &Room{name: roomName, visited: false, connections: []string{}}
			if parts[0] == startRoom {
				startRoom = roomName
			} else if parts[0] == endRoom {
				endRoom = roomName
			}
			continue
		}

		// Koridorları oku
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

	if err := scanner.Err(); err != nil {
		fmt.Println("Hata:", err)
	}

	// Verileri ekrana bas
	fmt.Println("Total Ants:", totalAnts)
	fmt.Println("Start Room:", startRoom)
	fmt.Println("End Room:", endRoom)
	fmt.Println("Rooms:", rooms)

	// Yolları bul
	paths := findPaths(rooms, startRoom, endRoom)
	fmt.Println("Bulunan Yollar:", paths)

	// Karıncaları yollar arasında dağıtma
	antPaths := assignAntsToPaths(totalAnts, paths)
	for ant, path := range antPaths {
		fmt.Printf("Karınca %d, Yol: %v\n", ant, path)
	}

	// Karıncaları hareket ettir ve adım adım ilerlet
	// Bu kısım da algoritmanıza göre tamamlanacaktır
}

func findPaths(rooms map[string]*Room, start string, end string) [][]string {
	var paths [][]string
	var queue [][]string

	queue = append(queue, []string{start})

	for len(queue) > 0 {
		path := queue[0]
		queue = queue[1:]

		lastRoom := path[len(path)-1]

		if lastRoom == end {
			paths = append(paths, path)
			continue
		}

		for _, connection := range rooms[lastRoom].connections {
			if !contains(path, connection) {
				newPath := append([]string{}, path...)
				newPath = append(newPath, connection)
				queue = append(queue, newPath)
			}
		}
	}

	return paths
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func assignAntsToPaths(totalAnts int, paths [][]string) map[int][]string {
	antPaths := make(map[int][]string)
	pathIndex := 0

	for ant := 1; ant <= totalAnts; ant++ {
		antPaths[ant] = paths[pathIndex]
		pathIndex++
		if pathIndex >= len(paths) {
			pathIndex = 0
		}
	}

	return antPaths
}
