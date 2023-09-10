package main

import (
	"fmt"
	"os"

	"github.com/emirpasic/gods/lists/arraylist"
)

var debug bool = true

// go run . {fileToRead} {DistanceToTravel}
func main() {
	//var workers int = 2  // Ammount of workers to read the data
	//var distance int = 0 // Distance to travel
	var File string = "Unknown"

	list := arraylist.New()
	list.Add("test")
	fmt.Println(list.Get(0))
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
		fmt.Println("No arguments provided")
		os.Exit(1)
	}

	// Initialize the data monitor
	dm := NewDataMonitor()

	var cars Cars = readFile(File)

	for i := 0; i < len(cars.Cars); i++ {
		dm.addDataItem(cars.Cars[i])
	}

	// Initialize the result monitor
	rm := NewResultMonitor()

	for i := len(dm.data) - 1; i >= 0; i-- {
		rm.addResultItem(dm.removeDataItem(0))
	}

	fmt.Println("Result:")
	result := rm.getResultItems()
	for i := 0; i < len(result); i++ {
		fmt.Println(result[i].Name)
	}
	fmt.Println("Count: ", len(result))
}
