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

	fmt.Println("Successfully Opened", filetoRead)
	// defer te closing of our jsonFile so that we can parse it later on
	defer inputFile.Close()

	// read our opened jsonFile as a byte array.
	byteValue, _ := io.ReadAll(inputFile)

	// initialize Cars array
	var cars Cars

	// unmarshal byteArray which contains
	// jsonFile's content into 'cars' which defined above
	json.Unmarshal(byteValue, &cars)

	fmt.Println("DEBUG | Cars found: " + strconv.Itoa(len(cars.Cars)))

	return cars
}

func writeFile(fileToWrite string, cars []Car) {
	// create a file for writing
	file, err := os.Create(fileToWrite)

	// if os.Create returns an error then handle it
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Successfully Created", fileToWrite)
	// defer te closing of our jsonFile so that we can parse it later on
	defer file.Close()

	carList := Cars{Cars: cars}

	// write our opened jsonFile as a byte array.
	byteValue, err := json.MarshalIndent(carList, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}

	// write the byteArray to our file
	_, err = file.Write(byteValue)
	if err != nil {
		fmt.Println(err)
		return
	}
}
