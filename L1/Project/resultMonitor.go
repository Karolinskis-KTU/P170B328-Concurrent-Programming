package main

import (
	"sort"
	"sync"
)

type ResultMonitor struct {
	sync.Mutex // Embed Mutex to synchronize access to the data (No race conditions)
	data       [30]Car
	count      int
}

func NewResultMonitor() *ResultMonitor {
	return &ResultMonitor{
		data:  [30]Car{},
		count: 0,
	}
}

// Adds a Car value to the ResultMonitor's data slice in sorted order by name.
// If the data slice is full, the program will stop and throw an error.
func (rm *ResultMonitor) addItemSorted(value Car) {
	rm.Lock()
	defer rm.Unlock()

	if rm.count == len(rm.data) {
		// Array is full, stop the program and throw an error
		panic("The result monitor is full.")
	}

	rm.data[rm.count] = value
	rm.count++

	// Sort the data slice by name
	sort.Slice((rm.data[:rm.count]), func(i, j int) bool {
		return rm.data[i].Name < rm.data[j].Name
	})
}

// Returns a copy of the ResultMonitor's data slice, sorted by name.
// The method creates a copy of te data slice to avoid concurrent access issues.
func (rm *ResultMonitor) getResultItems() []Car {
	rm.Lock()
	defer rm.Unlock()

	// Create a copy of the data slice to avoid concurrent access issues
	dataCopy := make([]Car, rm.count)
	copy(dataCopy, rm.data[:rm.count])

	return dataCopy
}
