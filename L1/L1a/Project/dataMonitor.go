package main

import (
	"sync"
)

type DataMonitor struct {
	data       []Car
	maxSize    int
	count      int
	cond       *sync.Cond
	isSignaled bool // Indicates if the data monitor has been signaled to stop
	sync.Mutex      // Embed Mutex to synchronize access to the data (No race conditions)
}

func NewDataMonitor(size int) *DataMonitor {
	dm := &DataMonitor{
		data:    make([]Car, 0, size),
		count:   0,
		maxSize: size,
	}
	dm.cond = sync.NewCond(&dm.Mutex)
	return dm
}

// Adds a Car value to the data monitor
// If the data monitor is full, the method will wait until an item is removed
func (dm *DataMonitor) addDataItem(value Car) {
	dm.Lock()
	defer dm.Unlock()

	if len(dm.data) >= dm.maxSize {
		// Array is full, wait until an item is removed
		dm.cond.Wait()
	}

	dm.data = append(dm.data, value)
	dm.count++

	dm.cond.Broadcast()
}

// Removes the last item from the data monitor
// If the data monitor is empty, the method will wait until an item is added
// If the data monitor is signaled to stop and there are no more data items in the data monitor, the method will return an empty Car value
func (dm *DataMonitor) removeDataItem() Car {
	dm.Lock()
	defer dm.Unlock()

	for len(dm.data) == 0 && !dm.isSignaled {
		dm.cond.Wait()
	}

	if dm.isSignaled && len(dm.data) == 0 {
		// The data monitor has been signaled to stop and there are no more data items in the data monitor
		return Car{}
	}

	removedItem := dm.data[len(dm.data)-1]
	dm.data = dm.data[:len(dm.data)-1]
	dm.count--

	dm.cond.Broadcast() // Broadcast that an item has been removed
	return removedItem
}

// Signals the data monitor to stop
func (dm *DataMonitor) signalStop() {
	dm.Lock()
	defer dm.Unlock()

	dm.isSignaled = true
	dm.cond.Broadcast()
}
