package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
)

var dm = NewDataMonitor()
var rm = NewResultMonitor()

// go run . {fileToRead} {DistanceToTravel}
func main() {
	var workers int = 10 // Ammount of workers to read the data
	//var distance int = 0 // Distance to travel
	var File string = "Unknown"

	if len(os.Args) > 1 {
		if os.Args[1] == "1" {
			// All cars fit the criteria
			File = "../Data/IFF-1-1_PaulaviciusK_L1_dat_1.json"
		} else if os.Args[1] == "2" {
			// Only some of the cars fit the criteria
			File = "../Data/IFF-1-1_PaulaviciusK_L1_dat_2.json"
		} else if os.Args[1] == "3" {
			// None of the cars fit the criteria
			File = "../Data/IFF-1-1_PaulaviciusK_L1_dat_3.json"
		}
	} else {
		fmt.Println("ERROR | No arguments provided")
		os.Exit(1)
	}

	// Initialize the data monitor
	//dm := NewDataMonitor()

	var cars Cars = readFile(File)

	for i := 0; i < len(cars.Cars); i++ {
		dm.addDataItem(cars.Cars[i])
	}

	// Generate worker threads
	fmt.Println("Starting ", workers, " workers...")
	var name string = "Worker"
	var waitGroup = sync.WaitGroup{}
	waitGroup.Add(workers)
	for i := 0; i < workers; i++ {
		var workerName string = name + strconv.Itoa(i)
		go execute(workerName, &waitGroup)
	}

	// Wait for all workers to finish
	waitGroup.Wait()

	fmt.Println("Result:")
	result := rm.getResultItems()
	fmt.Println("Name | Fuel Efficiency | Fuel Tank Size")
	for i := 0; i < len(result); i++ {
		fmt.Println(result[i].Name, "|", result[i].FuelEfficiency, "|", result[i].FuelTankSize)
	}
	fmt.Println("Count: ", len(result))

	writeFile("result.json", result)
}

func execute(name string, group *sync.WaitGroup) {
	//defer group.Done() // defer to the end of the function regardless if there is an error or not
	fmt.Println(name, "| Starting...")
	var loop int
	for {
		item := dm.removeDataItem(0)
		if item.Name == "" {
			break
		}

		rm.addItemSorted(item)

		loop++
	}
	/*
		var item = dm.removeDataItem(0)
		// Check if the returned item is empty
		for item.Name != "" {
			rm.addItemSorted(item)
			item = dm.removeDataItem(0)
			loop++
		}
	*/
	fmt.Println(name, "| Finished | Loops: ", loop)

	if group != nil {
		group.Done()
	}
}
