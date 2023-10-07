package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"sync"
)

// go run . {fileToRead}
func main() {
	var inputFile string = ""
	var outputFile string = "../Data/IFF-1-1_PaulaviciusK_L1_res.txt"

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "1":
			inputFile = "../Data/IFF-1-1_PaulaviciusK_L1_dat_1.json"
		case "2":
			inputFile = "../Data/IFF-1-1_PaulaviciusK_L1_dat_2.json"
		case "3":
			inputFile = "../Data/IFF-1-1_PaulaviciusK_L1_dat_3.json"
		}
	} else {
		fmt.Println("ERROR | No arguments provided")
		os.Exit(1)
	}

	// Read the data file and add te data items to te data monitor
	cars := readFile(inputFile)

	carCount := len(cars.Cars)
	workersCount := int(math.Max(2, float64(carCount/4)))

	dm := NewDataMonitor(carCount / 4)
	rm := NewResultMonitor(carCount)

	// Start the workers
	var waitGroup sync.WaitGroup
	for i := 0; i < workersCount; i++ {
		waitGroup.Add(1)
		go execute(fmt.Sprintf("Worker %d", i), dm, rm, &waitGroup)
	}
	println("Starting " + strconv.Itoa(workersCount) + " workers")
	for i := 0; i < carCount; i++ {
		dm.addDataItem(cars.Cars[i])
	}
	dm.signalStop()

	waitGroup.Wait()

	if _, err := os.Stat(outputFile); err == nil {
		os.Remove(outputFile)
	}

	result := rm.getResultItems()
	printData(cars.Cars, "Data")
	printData(result, "Result")
	writeFile(outputFile, cars.Cars, "Data")
	writeFile(outputFile, result, "Result")
}

func execute(name string, dm *DataMonitor, rm *ResultMonitor, group *sync.WaitGroup) {
	defer group.Done()

	for {
		if len(dm.data) == 0 {
			if dm.isSignaled {
				return
			}
		} else {
			item := dm.removeDataItem()
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
