package main

import (
	"fmt"
	"hash/fnv"
	"strings"
)

// Cars struct which contains
// an array of cars
type Cars struct {
	Cars []Car `json:"cars"`
}

type Car struct {
	Name           string  `json:"name"`
	FuelTankSize   int     `json:"fuel_tank_size"`
	FuelEfficiency float64 `json:"fuel_efficiency"`
}

type CarComputed struct {
	Car      Car `json:"car"`
	HashCode int `json:"hash_code"`
}

func (c Car) hashCode() int {
	data := strings.Join([]string{c.Name, fmt.Sprint(c.FuelTankSize), fmt.Sprint(int(c.FuelEfficiency))}, "")

	h := fnv.New32a()
	h.Write([]byte(data))
	hashCode := int(h.Sum32())

	return hashCode
}
