package main

import (
	"sort"
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

func (rm *ResultMonitor) addItemSorted(value Car) {
	rm.Lock()
	defer rm.Unlock()

	var carComputed CarComputed
	carComputed.Car = value
	carComputed.HashCode = value.hashCode()

	// Check if the sum of the digits of the hash code is even
	temp := carComputed.HashCode
	sum := 0
	for temp != 0 {
		sum += temp % 10
		temp /= 10
	}

	// If the sum is even, add the item to the result list
	if sum%2 == 0 {

		rm.data = append(rm.data, value)

		sort.Slice(rm.data, func(i, j int) bool {
			return rm.data[i].Name < rm.data[j].Name
		})
		//fmt.Println("Added item: ", value.Name)
	} else {
		//fmt.Println("Ignored item: ", value.Name)
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
