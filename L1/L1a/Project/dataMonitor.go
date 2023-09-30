package main

import (
	"sync"
)

type DataMonitor struct {
	sync.Mutex // Embed Mutex to synchronize access to the data (No race conditions)
	data       [10]Car
	count      int
	cond       *sync.Cond
}

func NewDataMonitor() *DataMonitor {
	dm := &DataMonitor{
		data:  [10]Car{},
		count: 0,
	}
	dm.cond = sync.NewCond(&dm.Mutex)
	return dm
}

// Adds a Car value to the DataMonitor's data slice.
// If the data slice is full, the program will wait until an item is removed.
// If the added item is the first item in the data slice, the method will signal that an item has been added.
func (dm *DataMonitor) addDataItem(value Car) {
	dm.Lock()
	defer dm.Unlock()

	for dm.count == len(dm.data) {
		// Array is full, wait for an item to be removed
		dm.Unlock()
		dm.Lock()
	}

	dm.data[dm.count] = value
	dm.count++

	if dm.count == 1 {
		// Signal that an item has been added
		dm.cond.Signal()
	}
}

// Removes a Car value from the DataMonitor's data slice at the given index.
// If an item is removed and the data slice is no longer full, the method will signal that an item has been removed.
func (dm *DataMonitor) removeDataItem(index int) Car {
	dm.Lock()
	defer func() {
		dm.Unlock()
	}()

	if index >= 0 && index < dm.count {
		removedItem := dm.data[index]
		copy(dm.data[index:], dm.data[index+1:dm.count])
		dm.count--
		//fmt.Println("Removed item: ", removedItem.Name)

		if dm.count == len(dm.data)-1 {
			// Signal that an item has been removed
			dm.cond.Signal()
		}

		return removedItem
	}

	return Car{Name: "", FuelTankSize: 0, FuelEfficiency: 0}
}
