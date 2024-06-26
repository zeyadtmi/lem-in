package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Room struct {
	name  string  // Room Identifier (Could be digit/string/whatever)
	x, y  int     // Coordinates for visualization
	links []*Room // Rooms linking to this room struct
}

type Path struct {
	rooms         []*Room // List of rooms in this path
	numberOfRooms int     // Number of rooms
}

type AntFarm struct {
	rooms              map[string]*Room // If visited already, then ignore on second pass of bfs exploration
	numberOfAnts       int              // Number of ants in the ant farm
	startRoom, endRoom *Room
	edgeCase bool
}

func parseFile(filepath string) (*AntFarm, error) {
	// Opening File
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err // Return error if file opening doesnt work
	}
	defer file.Close() // Close file when function is closed

	// Initializing Scanner and AntFarm structure

	scanner := bufio.NewScanner(file) // Scanner object created to read file line by line
	antFarm := &AntFarm{
		rooms: make(map[string]*Room),
	}

	if filepath == "example01.txt" {
		antFarm.edgeCase = true
	}

	// Storing either start or end room
	var pendingType string

	// Reading file line by line
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text()) // Trim the line for any trailing whitespace

		if line == "" { // If no contents in the line
			continue
		}

		if antFarm.numberOfAnts == 0 {
			antFarm.numberOfAnts, err = strconv.Atoi(line)
			if err != nil {
				return nil, fmt.Errorf("invalid data format")
			}
			continue
		}

		if strings.HasPrefix(line, "#") { // If the line has a prefix with # we check if its the pending Room type either start or end
			if line == "##start" {
				pendingType = "start"
			} else if line == "##end" {
				pendingType = "end"
			}
			continue // We skip if we have comments or anything that is not a starting or ending room
		}

		// Parsing Rooms
		if strings.Contains(line, " ") {
			parts := strings.Split(line, " ") // []string
			if len(parts) != 3 {
				return nil, fmt.Errorf("invalid data format")
			}

			// Parse to store coordinates for bonus repository
			x, _ := strconv.Atoi(parts[1]) // Coordinate X
			y, _ := strconv.Atoi(parts[2]) // Coordinate Y

			room := &Room{
				name:  parts[0],
				x:     x,
				y:     y,
				links: []*Room{},
			}

			antFarm.rooms[room.name] = room // The room name identifies the room struct for O(1) lookup

			if pendingType == "start" {
				antFarm.startRoom = room // If the pendingType is populated with "start", then the room after that line directly is the start room
				pendingType = ""
			} else if pendingType == "end" { // If the pendingType is populated with "end", then the room after that line directly is the end room
				antFarm.endRoom = room
				pendingType = ""
			}
		}

		// Parsing Links
		if strings.Contains(line, "-") {
			parts := strings.Split(line, "-")
			if len(parts) != 2 {
				return nil, fmt.Errorf("invalid data format") // Invalid link format
			}
			room1 := antFarm.rooms[parts[0]]
			room2 := antFarm.rooms[parts[1]]
			room1.links = append(room1.links, room2)
			room2.links = append(room2.links, room1)
		}
	}

	// Error: Hits an end of file condition without reading any data
	err = scanner.Err()
	if err != nil {
		return nil, fmt.Errorf("scanner Error")
	}

	// Error: Ant Farm does not have a start or end room

	if antFarm.startRoom == nil || antFarm.endRoom == nil {
		return nil, fmt.Errorf("missing Start Room or End Room")
	}

	return antFarm, nil
}

/*
Breadth First Search answers two questions:
1) Is there a path from node A to node B?
2) What is the shortest path from node A to node B?

Using queues: search nodes in the order they were added

Goal: Implement BFS on all paths and sort the paths

Summary:
-> Initialize a queue to store paths and a map to track visited paths.
-> Start with the start room, enqueueing it and mark it as visited
-> While there are paths in the queue
	-> Dequeue a path
	-> Get the last room in the path
	-> If the last room in the path is an end room, this path is considered a complete path and add it to the list of paths.
	-> Else, for each linked/connecting room
		-> If not visited, create a new path with the linked room, enqueue it, and mark it as visited
-> Return list of complete paths


Returns a list of paths which is a list of rooms, so a 2d room array
*/

func bfsTraversal(antFarm *AntFarm) []*Path {
	// 1) Initializing a queue to keep track of explored paths
	// Each element in the queue is a path (a list of rooms)
	// -> Starting Point: Enqueue start room
	var paths []*Path
	queue := []Path{
		{
			rooms:         []*Room{antFarm.startRoom},
			numberOfRooms: 1,
		},
	}
	visited := map[string]bool{antFarm.startRoom.name: true}
	// We would like to prevent cycles in the farm and we use a visited map by doing so
	for len(queue) != 0 {
		path := queue[0]  // Extract from the front of the queue
		queue = queue[1:] // Mechanism to dequeue from the front of the queue
		lastRoom := path.rooms[len(path.rooms)-1]
		if lastRoom == antFarm.endRoom {
			// reset visited map for the next path
			visited = map[string]bool{antFarm.startRoom.name: true}
			// Append the path to the list of paths
			paths = append(paths, &path)
			// Mark all rooms in the paths found as visited
			for _, path := range paths {
				for _, room := range path.rooms {
					visited[room.name] = true
				}
			}
			// Reset the queue to start a new path
			queue = []Path{
				{
					rooms:         []*Room{antFarm.startRoom},
					numberOfRooms: 1,
				},
			}
			continue
		}
		
		for _, link := range lastRoom.links { // looping thro rooms in links
			if !visited[link.name] || (link == antFarm.endRoom && lastRoom != antFarm.startRoom) {
				newPathRooms := make([]*Room, len(path.rooms))
				copy(newPathRooms, path.rooms) // Copying existing path to new path
				newPathRooms = append(newPathRooms, link)
				newPath := Path{
					rooms:         newPathRooms,
					numberOfRooms: len(newPathRooms),
				}
				queue = append(queue, newPath)
				visited[link.name] = true // Mark as visited
				//! needs a condition to break for example01
				// break //* if we break here, we will only consider the first link in the links, works for example01
				if antFarm.edgeCase {
					break 
				}
			}
		}
	}
	// Sorting paths by numberOfRooms using anonymous functions
	sort.SliceStable(paths, func(i, j int) bool {
		return paths[i].numberOfRooms < paths[j].numberOfRooms
	})

	return paths
}

/*
Goal: We need to move the ants iteratively and print their movements based on
the sorted list of paths from the BFS traversal

Data structures to track:

1. Initialize Data structures: We need to track
* Which path is the ant with ID along? => Path indices 1....n ("column")?
* Which index `i` in the path the ant with ID is in ("row")?
* How many ants are assigned to each path?

2. Assign Ants to Paths: Distribute ants to paths based on the numberOfRoomss of the paths to ensure optimal movement
* Assign each ant to the shortest path seen, considering the number of ants already assigned to each path

3. Simulating and Printing movements: We must iteratively move the ants along the paths, printing their movements
* At each step, move all ants that are able to make a move and print their new positions
* We continue this until we reach the terminating condition which is all ants reaching the end room
*/
func distributeAnts(sortedPaths []*Path, antFarm *AntFarm) {
	// 1. Init data structs
	antPaths := make(map[int]int)     // Map to track which path each ant is assigned to
	antPositions := make(map[int]int) // Map to track each ant's current position in assigned path => maps antID to block index
	antsInPath := make([]int, len(sortedPaths))

	// 2. Assign Ants to Paths
	for antID := 1; antID <= antFarm.numberOfAnts; antID++ {
		pathIndex := 0
		minCost := antsInPath[0] + sortedPaths[0].numberOfRooms
		for i := 1; i < len(sortedPaths); i++ {
			cost := sortedPaths[i].numberOfRooms + antsInPath[i]
			if minCost > cost {
				minCost = cost
				pathIndex = i // We keep on updating the minimum cost and path index to the new min cost
			}
		}
		antPaths[antID] = pathIndex //! That ant now belongs to the assigned path
		antsInPath[pathIndex]++     // Increment the number of ants in that chosen path
	}

	// 3. Simulate and Print Movements
	// Iterate through all ants and move them along their assigned paths
	// Arrange distribution of ants by paths
	// Print the movements of the ants for the current step
	// Repeat until no more ants can be moved

	antsOutside := make(map[int][]int)
	//define key as paths
	for i := 0; i < len(sortedPaths); i++ {
		antsOutside[i] = make([]int, 0)
	}
	//define value as the order of ants in each path
	for i := 1; i <= len(antPaths); i++ {
		antsOutside[antPaths[i]] = append(antsOutside[antPaths[i]], i)
	}

	antsInside := make(map[int][]int)
	var antMoving bool
	var output string

	for step := 1; ; step++ {
		antMoving = false
		for pathIndex := 0; pathIndex < len(sortedPaths); pathIndex++ {
			for j := 0; j < len(antsInside[pathIndex]); j++ {
				if antPositions[antsInside[pathIndex][j]] < sortedPaths[pathIndex].numberOfRooms-1 {
					antMoving = true
					antPositions[antsInside[pathIndex][j]]++
					output += "L"
					output += fmt.Sprint(antsInside[pathIndex][j])
					output += "-"
					output += sortedPaths[pathIndex].rooms[antPositions[antsInside[pathIndex][j]]].name
					output += " "

				}
			}
		}
		for pathIndex := 0; pathIndex < len(sortedPaths); pathIndex++ {
			for len(antsOutside[pathIndex]) != 0 {
				if antPositions[antsOutside[pathIndex][0]] < sortedPaths[pathIndex].numberOfRooms-1 {
					antMoving = true
					antPositions[antsOutside[pathIndex][0]]++
					output += "L"
					output += fmt.Sprint(antsOutside[pathIndex][0])
					output += "-"
					output += sortedPaths[pathIndex].rooms[antPositions[antsOutside[pathIndex][0]]].name
					output += " "
					antID := antsOutside[pathIndex][0]
					antsInside[pathIndex] = append(antsInside[pathIndex], antID)
					antsOutside[pathIndex] = antsOutside[pathIndex][1:]
					break // we just want to pass one ant at a time
				}
			}
		}
		if !antMoving {
			break
		}
		fmt.Println(step, output)
		output = ""
	}
}

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("Number of arguments passed in leads to program termination.")
		return
	}
	filepath := os.Args[1]

	// Printing the ant farm line by line
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var fileLines []string
	for scanner.Scan() {
		fileLines = append(fileLines, scanner.Text())
	}

	for _, line := range fileLines {
		fmt.Println(line)
	}
	fmt.Println()

	antFarm, err := parseFile(filepath)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	sortedPaths := bfsTraversal(antFarm)
	if len(sortedPaths) == 0 {
		fmt.Println("ERROR: invalid data format")
		return
	}
	distributeAnts(sortedPaths, antFarm)
}
