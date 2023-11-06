package main

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

// Car
type Car struct {
	Name           string  `json:"name"`
	FuelTankSize   int     `json:"fuel_tank_size"`
	FuelEfficiency float64 `json:"fuel_efficiency"`
	Fibbonaci      int     // Fibbonaci number
}

type Cars struct {
	Cars []Car `json:"cars"`
}

func (c Car) hashCode() int {
	data := strings.Join([]string{c.Name, fmt.Sprint(c.FuelTankSize), fmt.Sprint(int(c.FuelEfficiency))}, "")

	h := fnv.New32a()
	h.Write([]byte(data))
	hashCode := int(h.Sum32())

	return hashCode
}

// Settings
const workerCount = 6
const DataProcessSize = 10
const DataFileSize = 25

// go run . {fileToRead}
func main() {
	var inputFile string = ""
	var outputFile string = ""

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "1":
			// All fit the criteria
			inputFile = "../Data/IFF-1-1_PaulaviciusK_L2_dat_1.json"
			outputFile = "../Data/IFF-1-1_PaulaviciusK_L2_res_1.txt"
		case "2":
			// Some fit the criteria
			inputFile = "../Data/IFF-1-1_PaulaviciusK_L2_dat_2.json"
			outputFile = "../Data/IFF-1-1_PaulaviciusK_L2_res_2.txt"
		case "3":
			// None fit the criteria
			inputFile = "../Data/IFF-1-1_PaulaviciusK_L2_dat_3.json"
			outputFile = "../Data/IFF-1-1_PaulaviciusK_L2_res_3.txt"
		default:
			fmt.Println("ERROR | Invalid argument provided")
			os.Exit(1)
		}
	} else {
		fmt.Println("ERROR | No arguments provided")
		os.Exit(1)
	}

	// Read data from the file
	var cars Cars = ReadFile(inputFile)

	// Channels
	dataInputChannel := make(chan Car)     // Pradiniu duomenu siuntimas i duomenu kanala is pagrindines gijos i duomenu gijos
	dataOutputChannel := make(chan Car)    // Duoemnu siuntimas i darbuotoju kanala is duoemnu gijos i darbuotoju gija
	dataRemoveRequest := make(chan byte)   // Prasymas istrinti duomenis is duomenu gijos
	dataAdditionRequest := make(chan byte) // Prasymas prideti duomenis i duomenu gija

	carsOutputChannel := make(chan Car)  // Duomenys praeje filtra siunciami i rezultatu kanala is darbuotoju gijos
	carsToMainChannel := make(chan byte) // Prasymas baigti darba is darbuotoju gijos i pagrindine gija

	ResultsOutputChannel := make(chan Car) // Rezultatu kanalas is rezultatu gijos i pagrindine gija

	// Worker threads
	for i := 0; i < workerCount; i++ {
		go func() {
			for {
				dataRemoveRequest <- '+'   // Send a data removal request
				Car := <-dataOutputChannel // Receive car data from the data thread
				if Car.Name == "<Error>" {
					break
				}

				hashCode := Car.hashCode()
				fib := closestFibonacci(hashCode)

				if sumIsEven(fib) {
					carsOutputChannel <- Car // Send car data to the results thread
				}
			}
			carsToMainChannel <- '+' // Send a worker completion request
		}()
	}

	// Data thread
	go func() {
		defer close(dataOutputChannel)
		var CarClone [DataProcessSize]Car
		var Constants = [3]int{0, 0, 0} // Start, End, Size

		AddCar := func() {
			// Receive car data from the main thread
			// Update the "End" and "Size" constants
			CarClone[Constants[1]] = <-dataInputChannel
			Constants[1] = (Constants[1] + 1) % DataProcessSize
			Constants[2]++
		}

		RemoveCar := func() {
			// Send car data to the worker threads
			// Update the "Start" and "Size" constants
			dataOutputChannel <- CarClone[Constants[0]]
			Constants[0] = (Constants[0] + 1) % DataProcessSize
			Constants[2]--
		}
		for {
			if Constants[2] > 0 && Constants[2] < DataProcessSize { // We have data, and the clone is not full
				select {
				case <-dataAdditionRequest: // If a data addition request is received, add a car
					AddCar()
				case <-dataRemoveRequest: // If a data removal request is received, remove a car
					RemoveCar()
				}
			} else if Constants[2] == 0 { // Car clone is empty
				Message := <-dataAdditionRequest // wait for a data addition request
				if Message == '-' {
					break
				}
				AddCar()
			} else {
				<-dataRemoveRequest // wait for a data removal request
				RemoveCar()
			}
		}
	}()

	// Results thread
	go func() {
		defer close(ResultsOutputChannel)
		var Results [DataFileSize]Car
		Count := 0
		for Car := range carsOutputChannel { // Continuously receive car data from the worker threads
			i := Count
			for i > 0 && ((Results[i-1].FuelTankSize == Car.FuelTankSize && int(Results[i-1].FuelEfficiency) > int(Car.FuelEfficiency)) || Results[i-1].FuelTankSize > Car.FuelTankSize) {
				// Sort the car data
				Results[i] = Results[i-1]
				i--
			}
			Results[i] = Car
			Count++
		}
		for i := 0; i < Count; i++ {
			// Send the sorted car data to the main thread
			ResultsOutputChannel <- Results[i]
		}
	}()

	// Loop through the cars array in the 'cars' variable
	// Send the car data to the data thread
	for i := 0; i < len(cars.Cars); i++ {
		dataAdditionRequest <- '+'
		dataInputChannel <- cars.Cars[i]
		time.Sleep(300 * time.Millisecond)
	}

	// Create a termination signal for the worker threads
	for i := 0; i < workerCount; i++ {
		dataAdditionRequest <- '+'
		dataInputChannel <- Car{Name: "<Error>"}
		<-carsToMainChannel
	}

	close(carsToMainChannel)
	close(carsOutputChannel)

	dataAdditionRequest <- '-'

	close(dataInputChannel)
	close(dataAdditionRequest)
	close(dataRemoveRequest)

	var Results []Car
	for car := range ResultsOutputChannel {
		Results = append(Results, car)
	}

	// Print results
	// If output file is already created, delete it
	if _, err := os.Stat(outputFile); err == nil {
		os.Remove(outputFile)
	}
	WriteFile(outputFile, cars.Cars, "Initial Data")
	WriteFile(outputFile, Results, "Results")

	// PrintData(cars.Cars, "Initial Data")
	// PrintData(Results, "Results")
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

func sumIsEven(n int) bool {
	sum := 0
	for n > 0 {
		sum += n % 10
		n /= 10
	}

	return sum%2 == 0
}

func ReadFile(filetoRead string) Cars {
	inputFile, err := os.Open(filetoRead)

	// if os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("IO | Successfully opened", filetoRead)
	// defer te closing of our jsonFile so that we can parse it later on
	defer inputFile.Close()

	// read our opened jsonFile as a byte array.
	byteValue, _ := io.ReadAll(inputFile)

	// initialize Cars array
	var cars Cars

	// unmarshal byteArray which contains
	// jsonFile's content into 'cars' which defined above
	json.Unmarshal(byteValue, &cars)

	fmt.Println("IO | Cars found: " + strconv.Itoa(len(cars.Cars)))
	return cars
}

func WriteFile(fileToWrite string, cars []Car, header string) {
	// Create a file for writing
	file, err := os.OpenFile(fileToWrite, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)

	// If os.Create returns an error then handle it
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("IO | Successfully added to:", fileToWrite)
	// Defer the closing of our file so that we can close it later on
	defer file.Close()

	// Write header information to the file
	fmt.Fprintf(file, "%+30s\n", header)
	fmt.Fprintf(file, "%-5s | %-20s | %-17s | %-15s\n", "#", "Name", "Fuel Efficiency", "Fuel Tank Size")
	fmt.Fprintln(file, "------------------------------------------------------------------")

	// Write car data to the file
	for i, car := range cars {
		fmt.Fprintf(file, "%-5d | %-20s | %-17.2f | %-15d\n", i+1, car.Name, car.FuelEfficiency, car.FuelTankSize)
	}

	fmt.Fprintln(file, "Count: ", len(cars))
	fmt.Fprintln(file, "------------------------------------------------------------------")
}

func PrintData(cars []Car, header string) {
	fmt.Printf("%+30s \n", header)
	fmt.Printf("%-5s | %-20s | %-17s | %-15s\n", "#", "Name", "Fuel Efficiency", "Fuel Tank Size")
	fmt.Println("------------------------------------------------------------------")
	for i := 0; i < len(cars); i++ {
		fmt.Printf("%-5d | %-20s | %-17.2f | %-15d\n", i+1, cars[i].Name, cars[i].FuelEfficiency, cars[i].FuelTankSize)
	}
	fmt.Println("Count: ", len(cars))
	fmt.Println("------------------------------------------------------------------")
}
