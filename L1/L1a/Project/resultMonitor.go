package main

import (
	"sort"
	"sync"
)

type ResultMonitor struct {
	data       []Car
	count      int
	sync.Mutex // Embed Mutex to synchronize access to the data (No race conditions)
}

func NewResultMonitor(maxSize int) *ResultMonitor {
	return &ResultMonitor{
		data:  make([]Car, 0, maxSize),
		count: 0,
	}
}

// Adds a Car value to the result monitor
// If the result monitor is full, the method will stop the program and throw an error
func (rm *ResultMonitor) addItemSorted(value Car) {
	rm.Lock()
	defer rm.Unlock()

	if len(rm.data) == 0 {
		rm.data = append(rm.data, value)
	} else {
		i := sort.Search(len(rm.data), func(i int) bool { return rm.data[i].Name < value.Name })
		// If i is start index just adds it to beginning
		if i == 0 {
			rm.data = append([]Car{value}, rm.data...)
		} else {
			// Otherwise just add estate to position where it belongs
			rm.data = append(rm.data, Car{})
			copy(rm.data[i+1:], rm.data[i:])
			rm.data[i] = value
		}
	}
	rm.count++
}

// Returns a copy of the ResultMonitor's data slice, sorted by name.
// The method creates a copy of the data slice to avoid concurrent access issues.
func (rm *ResultMonitor) getResultItems() []Car {
	// Create a copy of the data slice to avoid concurrent access issues
	dataCopy := make([]Car, rm.count)
	copy(dataCopy, rm.data[:rm.count])

	return dataCopy
}
