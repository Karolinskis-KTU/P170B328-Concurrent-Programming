package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
)

func readFile(filetoRead string) Cars {
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

func writeFile(fileToWrite string, cars []Car, header string) {
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

func printData(cars []Car, header string) {
	fmt.Printf("%+30s \n", header)
	fmt.Printf("%-5s | %-20s | %-17s | %-15s\n", "#", "Name", "Fuel Efficiency", "Fuel Tank Size")
	fmt.Println("------------------------------------------------------------------")
	for i := 0; i < len(cars); i++ {
		fmt.Printf("%-5d | %-20s | %-17.2f | %-15d\n", i+1, cars[i].Name, cars[i].FuelEfficiency, cars[i].FuelTankSize)
	}
	fmt.Println("Count: ", len(cars))
	fmt.Println("------------------------------------------------------------------")
}
