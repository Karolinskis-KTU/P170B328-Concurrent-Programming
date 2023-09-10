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

	if debug {
		fmt.Println("Cars found: " + strconv.Itoa(len(cars.Cars)))
	}

	return cars
}
