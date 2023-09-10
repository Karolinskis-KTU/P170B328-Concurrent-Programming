package main

import (
	"fmt"
	"sync"
)

type ResultMonitor struct {
	sync.Mutex // Embed Mutex to synchronize access to the data (No race conditions)
	data       []Car
}

func NewResultMonitor() *ResultMonitor {
	return &ResultMonitor{
		data: []Car{},
	}
}

func (rm *ResultMonitor) addResultItem(value Car) {
	rm.Lock()
	defer rm.Unlock()

	var carComputed CarComputed
	carComputed.Car = value
	carComputed.HashCode = value.hashCode()

	temp := carComputed.HashCode
	sum := 0
	for temp != 0 {
		sum += temp % 10
		temp /= 10
	}

	if sum%2 == 0 {
		rm.data = append(rm.data, value)
		fmt.Println("Added item: ", value.Name)
	} else {
		fmt.Println("Ignored item: ", value.Name)
	}

}

func (rm *ResultMonitor) getResultItems() []Car {
	rm.Lock()
	defer rm.Unlock()

	// Create a copy of the data slice to avoid concurrent access issues
	dataCopy := make([]Car, len(rm.data))
	copy(dataCopy, rm.data)

	return dataCopy
}
