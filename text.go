package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type room struct {
	name     string
	coord    [2]int
	antCount int
}

type connection struct {
	from string
	to   string
}

type antFarm struct {
	antTotal int
	rooms    map[string]room
	links    map[string][]string
	start    string
	end      string
}

func newAntFarm() *antFarm {
	return &antFarm{
		rooms: make(map[string]room),
		links: make(map[string][]string),
	}
}

func (farm *antFarm) parseInput(inputLines []string) error {
	if len(inputLines) < 1 {
		return fmt.Errorf("ERROR: invalid data format, input is empty")
	}

	// Parse number of ants
	var err error
	farm.antTotal, err = strconv.Atoi(strings.TrimSpace(inputLines[0]))
	if err != nil || farm.antTotal <= 0 {
		return fmt.Errorf("ERROR: invalid data format, invalid number of ants")
	}

	// Parse rooms and links
	for i := 1; i < len(inputLines); i++ {
		line := strings.TrimSpace(inputLines[i])

		switch {
		case line == "##start":
			i++
			line = strings.TrimSpace(inputLines[i])
			if err := farm.parseRoom(line, true, false); err != nil {
				return err
			}
		case line == "##end":
			i++
			line = strings.TrimSpace(inputLines[i])
			if err := farm.parseRoom(line, false, true); err != nil {
				return err
			}
		case strings.Contains(line, " "):
			if err := farm.parseRoom(line, false, false); err != nil {
				return err
			}
		case strings.Contains(line, "-"):
			if err := farm.parseLink(line); err != nil {
				return err
			}
		default:
			return fmt.Errorf("ERROR: invalid data format")
		}
	}

	if farm.start == "" || farm.end == "" {
		return fmt.Errorf("ERROR: invalid data format, no start or end room found")
	}

	return nil
}

func (farm *antFarm) parseRoom(line string, isStart bool, isEnd bool) error {
	parts := strings.Fields(line)
	if len(parts) != 3 {
		return fmt.Errorf("ERROR: invalid data format, invalid room")
	}
	name := parts[0]
	x, err1 := strconv.Atoi(parts[1])
	y, err2 := strconv.Atoi(parts[2])
	if err1 != nil || err2 != nil {
		return fmt.Errorf("ERROR: invalid data format, invalid coordinates")
	}
	farm.rooms[name] = room{name: name, coord: [2]int{x, y}}
	if isStart {
		if farm.start != "" {
			return fmt.Errorf("ERROR: invalid data format, duplicate start room")
		}
		farm.start = name
	}
	if isEnd {
		if farm.end != "" {
			return fmt.Errorf("ERROR: invalid data format, duplicate end room")
		}
		farm.end = name
	}
	return nil
}

func (farm *antFarm) parseLink(line string) error {
	parts := strings.Split(line, "-")
	if len(parts) != 2 {
		return fmt.Errorf("ERROR: invalid data format, invalid link")
	}
	name1 := parts[0]
	name2 := parts[1]
	if _, ok := farm.rooms[name1]; !ok {
		return fmt.Errorf("ERROR: invalid data format, unknown room in link")
	}
	if _, ok := farm.rooms[name2]; !ok {
		return fmt.Errorf("ERROR: invalid data format, unknown room in link")
	}
	farm.links[name1] = append(farm.links[name1], name2)
	farm.links[name2] = append(farm.links[name2], name1)
	return nil
}

func bfs(farm *antFarm) []string {
	queue := []string{farm.start}
	paths := make(map[string][]string)
	paths[farm.start] = []string{farm.start}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for _, neighbor := range farm.links[current] {
			if _, ok := paths[neighbor]; !ok {
				paths[neighbor] = append(paths[current], neighbor)
				queue = append(queue, neighbor)
				if neighbor == farm.end {
					return paths[neighbor]
				}
			}
		}
	}

	return nil
}

func edmondsKarp(farm *antFarm) int {
	maxFlow := 0
	path := bfs(farm)

	for path != nil {
		maxFlow++
		for i := 0; i < len(path)-1; i++ {
			current := path[i]
			next := path[i+1]
			farm.links[current] = remove(farm.links[current], next)
			farm.links[next] = append(farm.links[next], current)
		}
		path = bfs(farm)
	}

	return maxFlow
}

func remove(slice []string, s string) []string {
	for i, v := range slice {
		if v == s {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func (farm *antFarm) simulateAnts() {
	// Her adımda bir önceki adımdaki karıncaların durumunu güncelleyeceğiz.
	for step := 0; farm.antTotal > 0; step++ {
		// Her adımda karıncaların hareketlerini güncellemek için bir döngü yapacağız.
		for antID := 1; antID <= farm.antTotal; antID++ {
			antRoom := farm.rooms[strconv.Itoa(antID)] // Karıncanın bulunduğu oda
			if antRoom.name == farm.end {              // Eğer karınca hedef odaya ulaştıysa, bir çıkış yapar.
				delete(farm.rooms, strconv.Itoa(antID))
				farm.antTotal--
			} else { // Değilse, karıncayı ileri doğru hareket ettiririz.
				nextRoom := farm.links[antRoom.name][0] // Karıncanın bir sonraki odası
				farm.rooms[nextRoom] = room{name: nextRoom, coord: farm.rooms[nextRoom].coord, antCount: farm.rooms[nextRoom].antCount + 1}
				antRoom.antCount-- // Önceki odadaki karınca sayısını azalt
				if antRoom.antCount == 0 {
					delete(farm.rooms, antRoom.name)
				}
			}
		}
		// Karınca hareketlerini bir adım sonrasını göstermek için yazdırın
		fmt.Printf("Step %d:\n", step+1)
		for roomName, r := range farm.rooms {
			fmt.Printf("Room: %s, Ant Count: %d\n", roomName, r.antCount)
		}
	}
}

func (farm *antFarm) printOutput() {
	// Sonuçları belirtilen formatta yazdırmak için bir döngü yapacağız.
	for step := 0; farm.antTotal > 0; step++ {
		var output strings.Builder
		for antID := 1; antID <= farm.antTotal; antID++ {
			antRoom := farm.rooms[strconv.Itoa(antID)] // Karıncanın bulunduğu oda
			if antRoom.name == farm.end {              // Eğer karınca hedef odaya ulaştıysa, bir çıkış yapar.
				delete(farm.rooms, strconv.Itoa(antID))
				farm.antTotal--
			} else { // Değilse, çıktıya karıncanın yeni pozisyonunu ekleyin
				nextRoom := farm.links[antRoom.name][0] // Karıncanın bir sonraki odası
				output.WriteString(fmt.Sprintf("L%d-%s ", antID, nextRoom))
			}
		}
		// Karıncaların adım sonrasını gösteren çıktıyı yazdırın
		fmt.Println(output.String())
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var inputLines []string

	for scanner.Scan() {
		inputLines = append(inputLines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
		return
	}

	farm := newAntFarm()
	if err := farm.parseInput(inputLines); err != nil {
		fmt.Println(err)
		return
	}

	farm.simulateAnts()
	farm.printOutput()
}
