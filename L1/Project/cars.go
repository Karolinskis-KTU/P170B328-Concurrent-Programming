package main

import (
	"fmt"
	"hash/fnv"
	"math/rand"
	"strings"
	"time"
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

	rand.Seed(time.Now().UnixNano())
	minSleep := 10
	maxSleep := 50
	sleepDuration := time.Duration(rand.Intn(maxSleep-minSleep+1)+minSleep) * time.Millisecond

	time.Sleep(sleepDuration)

	return hashCode
}
