package structs

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

type Ant struct {
	ID       int
	Position int      
	Path     []string 
}

type PathInfo struct {
	Path           []string 
	Ants           int      
	CompletionTime int      
}

type PathGroup struct {
	Paths       [][]string 
	MaxLength   int        
	MaxAnts     int       
	CurrentAnts []int     
}

type Room struct {
	Name string
	X, Y int
}

type AntFarm struct {
	Rooms      map[string]Room     
	Tunnels    map[string][]string 
	Start, End string        
}

type searchState struct {
	room string   
	path []string 
}

func distributeAntsOneGroup(totalAnts, maxAnts, maxLength, numPaths int) int {
	if totalAnts <= maxAnts {
		return maxLength
	}

	excessAnts := totalAnts - maxAnts
	additionalTurns := (excessAnts + numPaths - 1) / numPaths
	return maxLength + additionalTurns
}

func Distribute(totalAnts int, group *PathGroup) []Ant {
	ants := make([]Ant, totalAnts)

	group.CurrentAnts = make([]int, len(group.Paths))

	type pathSlot struct {
		index int 
		cost  int 
	}
	pathHeap := make([]pathSlot, len(group.Paths))

	for i, path := range group.Paths {
		pathHeap[i] = pathSlot{
			index: i,
			cost:  len(path),
		}
	}

	heapSize := len(pathHeap)

	heapify := func() {
		sort.Slice(pathHeap, func(i, j int) bool {
			return pathHeap[i].cost < pathHeap[j].cost
		})
	}

	for i := 0; i < totalAnts; i++ {
		heapify()
		minPath := pathHeap[0]

		ants[i] = Ant{
			ID:       i + 1,
			Path:     group.Paths[minPath.index],
			Position: 0,
		}

		group.CurrentAnts[minPath.index]++

		pathHeap[0].cost = len(group.Paths[minPath.index]) + group.CurrentAnts[minPath.index] - 1

		if heapSize > 1 && pathHeap[0].cost > pathHeap[1].cost {
			heapify()
		}
	}

	return ants
}

func BestGroup(ants int, groups []PathGroup) (int, int, [][]string) {
	minSteps := math.MaxInt32
	bestGroupIndex := -1
	var bestPaths [][]string

	for i, group := range groups {
		steps := distributeAntsOneGroup(ants, group.MaxAnts, group.MaxLength, len(group.Paths))
		if steps < minSteps {
			minSteps = steps
			bestGroupIndex = i
			bestPaths = group.Paths
		}
	}

	return bestGroupIndex, minSteps - 1, bestPaths
}

func (farm *AntFarm) FindAllPaths() [][]string {
	var allPaths [][]string

	startState := searchState{room: farm.Start, path: []string{farm.Start}}

	visitedStart := make(map[string]bool)
	visitedStart[farm.Start] = true

	var findPaths func(state searchState, visited map[string]bool)
	findPaths = func(state searchState, visited map[string]bool) {
		if state.room == farm.End {
			allPaths = append(allPaths, state.path)
			return
		}

		for _, nextRoom := range farm.Tunnels[state.room] {
			if !visited[nextRoom] {
				visitedCopy := make(map[string]bool)
				for key, value := range visited {
					visitedCopy[key] = value
				}
				visitedCopy[nextRoom] = true

				newPath := make([]string, len(state.path)+1)
				copy(newPath, state.path)
				newPath[len(newPath)-1] = nextRoom

				findPaths(searchState{room: nextRoom, path: newPath}, visitedCopy)
			}
		}
	}
	findPaths(startState, visitedStart)

	return allPaths
}

func NewAntFarm(startRoom, endRoom string, rooms map[string]Room, tunnels []string) *AntFarm {
	farm := &AntFarm{
		Rooms:   rooms,
		Tunnels: make(map[string][]string),
		Start:   startRoom,
		End:     endRoom,
	}

	for _, tunnel := range tunnels {
		parts := strings.Split(tunnel, "-")
		if len(parts) == 2 {
			farm.Tunnels[parts[0]] = append(farm.Tunnels[parts[0]], parts[1])
			farm.Tunnels[parts[1]] = append(farm.Tunnels[parts[1]], parts[0])
		}
	}

	return farm
}

func pathsIntersect(path1, path2 []string) bool {
	set := make(map[string]bool)
	for _, node := range path1[1 : len(path1)-1] {
		set[node] = true
	}
	for _, node := range path2[1 : len(path2)-1] {
		if set[node] {
			return true
		}
	}
	// fmt.Println(path1)
	// fmt.Println(path2)
	return false
}

func NonIntersecting(paths [][]string) []PathGroup {
	var groups []PathGroup
	n := len(paths)

	if n == 0 {
		return groups 
	}

	var findGroups func(currentIndex int, currentGroup []int)
	findGroups = func(currentIndex int, currentGroup []int) {
		if currentIndex == n {
			if len(currentGroup) == 0 {
				return
			}
			newGroup := make([][]string, len(currentGroup))
			maxLength := 0
			for i, idx := range currentGroup {
				newGroup[i] = paths[idx]
				if len(paths[idx]) > maxLength {
					maxLength = len(paths[idx])
				}
			}

			maxAnts := 0
			for _, path := range newGroup {
				maxAnts += maxLength - len(path) + 1
			}

			groups = append(groups, PathGroup{
				Paths: newGroup,
				MaxLength: maxLength,
				MaxAnts:   maxAnts,
			})
			return
		}

		conflict := false
		for _, idx := range currentGroup {
			if pathsIntersect(paths[currentIndex], paths[idx]) {
				conflict = true
				break
			}
		}

		if !conflict {
			findGroups(currentIndex+1, append(currentGroup, currentIndex))
		}

		findGroups(currentIndex+1, currentGroup)
	}

	findGroups(0, []int{})

	sort.Slice(groups, func(i, j int) bool {
		return groups[i].MaxLength < groups[j].MaxLength
	})

	return groups
}

func SimulateAnts(ants []Ant, startRoom, endRoom string) {
	turn := 0

	occupiedRooms := make(map[string]bool)
	occupiedTunnels := make(map[string]bool)

	allAntsFinished := false

	for !allAntsFinished {
		allAntsFinished = true
		moves := make([]string, 0)

		for k := range occupiedTunnels {
			delete(occupiedTunnels, k)
		}

		for i := range ants {
			ant := &ants[i]

			if ant.Position >= len(ant.Path)-1 {
				continue
			}

			nextPosition := ant.Position + 1
			currentRoom := ant.Path[ant.Position]
			nextRoom := ant.Path[nextPosition]
			tunnel := currentRoom + "-" + nextRoom

			if !occupiedRooms[nextRoom] && !occupiedTunnels[tunnel] {
				ant.Position = nextPosition
				moves = append(moves, fmt.Sprintf("L%d-%s", ant.ID, nextRoom))

				if nextRoom != endRoom {
					occupiedRooms[nextRoom] = true
				}
				occupiedTunnels[tunnel] = true
				allAntsFinished = false
			}
		}

		if len(moves) > 0 {
			fmt.Printf("Turn %d: %s\n", turn+1, strings.Join(moves, " "))

			// fmt.Printf("%s\n", strings.Join(moves, " "))
		}

		turn++
		if turn > 100000 { 
			fmt.Println("Breaking to prevent a potential infinite loop.")
			break
		}

		for room := range occupiedRooms {
			if room != startRoom && room != endRoom {
				delete(occupiedRooms, room)
			}
		}
	}
}
