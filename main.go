package main

import (
	"fmt"
	"lem-in/simulate/structs"
	"lem-in/data/file"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . <input_file>")
		return
	}
	filename := os.Args[1]

	numberOfAnts, startRoom, endRoom, rooms, tunnels, err := file.ParseInput(filename)
	if err != nil {
		fmt.Println("ERROR: invalid data format")
		fmt.Println(err)
		return
	}

	farm := structs.NewAntFarm(startRoom, endRoom, rooms, tunnels)
	paths := farm.FindAllPaths()

	groups := structs.NonIntersecting(paths)
	if len(groups) == 0 {
		fmt.Println("No viable path groups found.")
		return
	}

	bestGroupIndex, minSteps, bestPaths := structs.BestGroup(numberOfAnts, groups)
	bestGroup := groups[bestGroupIndex]

	fmt.Println("Ants:",numberOfAnts)
	fmt.Println("Start room:", startRoom)
	fmt.Println("End room:", endRoom)
	fmt.Printf("Best group chosen: Group %d with minimum steps required: %d\n", bestGroupIndex+1, minSteps)
	fmt.Println("Paths in the best group:")
	for i, path := range bestPaths {
		fmt.Printf("Path %d: %v\n", i+1, path)
	}

	ants := structs.Distribute(numberOfAnts, &bestGroup)

	structs.SimulateAnts(ants, farm.Start, farm.End)
}
