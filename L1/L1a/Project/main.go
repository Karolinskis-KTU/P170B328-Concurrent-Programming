package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
)

var dm = NewDataMonitor()
var rm = NewResultMonitor()

// go run . {fileToRead}
func main() {
	var workers int = 10 // Ammount of workers to read the data
	var File string = ""

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

	// Generate wroker threads
	fmt.Println("MAIN | Starting ", workers, " workers...")
	var name string = "Worker"
	var waitGroup = sync.WaitGroup{}
	waitGroup.Add(workers)

	// Create channels to signal the workers to start and stop processing data
	start := make(chan bool)
	read_end := make(chan bool)

	for i := 0; i < workers; i++ {
		var workerName string = name + strconv.Itoa(i)
		go execute(workerName, &waitGroup, start, read_end)
	}

	// Send the start signal to te worker threads
	close(start)

	// Read te data file and add te data items to te data monitor
	cars := readFile(File)
	for i := 0; i < len(cars.Cars); i++ {
		dm.addDataItem(cars.Cars[i])
	}

	// Close the stop channel to signal that no more data items will be added to the data monitor
	close(read_end)

	// Wait for all workers to finish
	waitGroup.Wait()

	fmt.Println("Result:")
	result := rm.getResultItems()
	fmt.Printf("%-15s | %-17s | %-15s\n", "Name", "Fuel Efficiency", "Fuel Tank Size")
	fmt.Println("--------------------------------------------------------------")
	for i := 0; i < len(result); i++ {
		fmt.Printf("%-15s | %-17.2f | %-15d\n", result[i].Name, result[i].FuelEfficiency, result[i].FuelTankSize)
	}
	fmt.Println("Count: ", len(result))

	writeFile("result.json", result)
}

func execute(name string, group *sync.WaitGroup, start chan bool, read_end chan bool) {
	defer group.Done()
	fmt.Println(name, "| Starting...")

	// Wait for the start signal
	<-start
	for {
		item := dm.removeDataItem(0)

		if item.Name == "" {
			select {
			case <-read_end:
				//fmt.Println(name, "Got signal to stop...")
				if item.Name == "" {
					fmt.Println(name, "| Stopping... No more data items in the data monitor.")
					return
				} else {
					//  If there are still data items in the data monitor, process them
					continue
				}
			default:
				//fmt.Println(name, "| The data monitor is empty. Waiting for more data items...")
				continue
			}
		}

		var carComputed CarComputed
		carComputed.Car = item
		carComputed.HashCode = item.hashCode()

		temp := carComputed.HashCode
		// Get the closest fibonacci number to the hash code
		fib := closestFibonacci(temp)
		// Check if the sum of the digits in the fibonacci number is even
		sum := 0
		for fib != 0 {
			sum += fib % 10
			fib /= 10
		}

		// If the sum is even, add the item to the result list
		if sum%2 == 0 {
			//fmt.Println(name, "| Adding item to result list:", sum, item.Name)
			rm.addItemSorted(item)
		} else {
			//fmt.Println(name, "| Item discarded:", sum, item.Name)
			continue
		}
	}
}

func closestFibonacci(n int) int {
	a := 0
	b := 1
	fib := 0

	for fib <= n {
		fib = a + b
		a = b
		b = fib
	}

	if n < fib-n {
		return a
	} else {
		return b
	}
}
