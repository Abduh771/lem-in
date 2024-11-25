package file

import (
	"bufio"
	"fmt"
	"lem-in/simulate/structs"
	"os"
	"strconv"
	"strings"
)

func ParseInput(filename string) (numberOfAnts int, startRoom, endRoom string, rooms map[string]structs.Room, tunnels []string, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, "", "", nil, nil, err
	}
	defer file.Close()

	rooms = make(map[string]structs.Room)
	var lines []string
	scanner := bufio.NewScanner(file)

	if scanner.Scan() {
		numberOfAnts, err = strconv.Atoi(scanner.Text())
		if err != nil {
			return 0, "", "", nil, nil, fmt.Errorf("ERROR: invalid number of ants")
		}
		if numberOfAnts < 1 {
			return 0, "", "", nil, nil, fmt.Errorf("ERROR: invalid number of ants. there should be at least one ant")
		}
	}

	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	var nextRoomIsStart, nextRoomIsEnd, startRoomAssigned, endRoomAssigned bool

	for _, line := range lines {
		if line == "##start" {
			if startRoomAssigned {
				return 0, "", "", nil, nil, fmt.Errorf("ERROR: multiple start rooms specified")
			}
			nextRoomIsStart = true
			continue
		} else if line == "##end" {
			if endRoomAssigned {
				return 0, "", "", nil, nil, fmt.Errorf("ERROR: multiple end rooms specified")
			}
			nextRoomIsEnd = true
			continue
		}

		if strings.Contains(line, "-") {
			tunnelParts := strings.Split(line, "-")
			if len(tunnelParts) != 2 {
				return 0, "", "", nil, nil, fmt.Errorf("ERROR: wrong data format in tunnel: %s", line)
			}
			tunnels = append(tunnels, line)
		} else if line != "" && !strings.HasPrefix(line, "#") {
			parts := strings.Split(line, " ")
			if len(parts) == 3 {
				name := parts[0]
				x, err := strconv.Atoi(parts[1])
				if err != nil {
					return 0, "", "", nil, nil, fmt.Errorf("ERROR: invalid data format, invalid X coordinate for room %s", name)
				}
				y, err := strconv.Atoi(parts[2])
				if err != nil {
					return 0, "", "", nil, nil, fmt.Errorf("ERROR: invalid data format, invalid Y coordinate for room %s", name)
				}
				rooms[name] = structs.Room{Name: name, X: x, Y: y}

				if nextRoomIsStart {
					startRoom = name
					startRoomAssigned = true
					nextRoomIsStart = false
				} else if nextRoomIsEnd {
					endRoom = name
					endRoomAssigned = true
					nextRoomIsEnd = false
				}
			}
		}
	}

	if err := validateAntFarm(startRoom, endRoom, rooms, tunnels); err != nil {
		return 0, "", "", nil, nil, err
	}

	if !startRoomAssigned || !endRoomAssigned {
		return 0, "", "", nil, nil, fmt.Errorf("ERROR: start or end room not specified")
	}

	return numberOfAnts, startRoom, endRoom, rooms, tunnels, scanner.Err()
}

func validateAntFarm(startRoom, endRoom string, rooms map[string]structs.Room, tunnels []string) error {
	if startRoom == "" || endRoom == "" {
		return fmt.Errorf("ERROR: invalid data format, no start or end room found")
	}

	for _, tunnel := range tunnels {
		parts := strings.Split(tunnel, "-")
		if parts[0] == parts[1] {
			return fmt.Errorf("ERROR: invalid data format, room links to itself: %s", tunnel)
		}
		if _, exists := rooms[parts[0]]; !exists {
			return fmt.Errorf("ERROR: invalid data format, link to unknown room: %s", tunnel)
		}
		if _, exists := rooms[parts[1]]; !exists {
			return fmt.Errorf("ERROR: invalid data format, link to unknown room: %s", tunnel)
		}
	}

	seenRooms := make(map[string]bool)
	for name := range rooms {
		if seenRooms[name] {
			return fmt.Errorf("ERROR: invalid data format, duplicated room: %s", name)
		}
		seenRooms[name] = true
	}

	return nil
}
