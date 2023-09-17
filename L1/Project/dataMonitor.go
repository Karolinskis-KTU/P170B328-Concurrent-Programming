package main

import (
	"sync"
)

type DataMonitor struct {
	sync.Mutex // Embed Mutex to synchronize access to the data (No race conditions)
	data       []Car
}

func NewDataMonitor() *DataMonitor {
	return &DataMonitor{
		data: []Car{},
	}
}

func (dm *DataMonitor) addDataItem(value Car) {
	dm.Lock()
	defer dm.Unlock()
	dm.data = append(dm.data, value)
	//fmt.Println("Added item: ", value.Name)
}

func (dm *DataMonitor) removeDataItem(index int) Car {
	dm.Lock()
	defer dm.Unlock()
	if index >= 0 && index < len(dm.data) {
		removedItem := dm.data[index]
		dm.data = append(dm.data[:index], dm.data[index+1:]...)
		//fmt.Println("Removed item: ", removedItem.Name)
		return removedItem
	}

	return Car{Name: "", FuelTankSize: 0, FuelEfficiency: 0}
}
