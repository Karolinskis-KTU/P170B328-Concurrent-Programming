package main

import (
	"fmt"
	"sync"
)

type DataMonitor struct {
	maxSize        int
	reCount        int
	letters        string
	writeConsonant bool
	mu             sync.Mutex
}

func NewDataMonitor() *DataMonitor {
	return &DataMonitor{
		letters:        "*",
		reCount:        0,
		writeConsonant: false,
	}
}

func (d *DataMonitor) GetResults() {
	if d.reCount > 2 && d.letters[d.reCount] == 'A' && d.letters[d.reCount-1] == 'A' && d.letters[d.reCount-2] == 'A' {
		d.writeConsonant = true
	} else {
		d.writeConsonant = false
	}
}

func (d *DataMonitor) AddToResults(letter rune) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.GetResults()

	if d.writeConsonant && letter != 'A' {
		d.reCount++
		d.letters += string(letter)
	} else if !d.writeConsonant {
		d.reCount++
		d.letters += string("A")
	}
}

func (d *DataMonitor) Print() {
	fmt.Println(d.letters, len(d.letters))
}

func (d *DataMonitor) enoughLetters() bool {
	aCount := 0
	bCount := 0
	cCount := 0
	for _, letter := range d.letters {
		switch letter {
		case 'A':
			aCount++
		case 'B':
			bCount++
		case 'C':
			cCount++
		}
	}
	return aCount >= 15 || bCount >= 15 || cCount >= 15
}

func (d *DataMonitor) worker(letter rune, wg *sync.WaitGroup) {
	defer wg.Done()
	for !d.enoughLetters() {
		d.AddToResults(letter)
	}
}

func main() {
	monitor := NewDataMonitor()
	threadNo := 3

	var wg sync.WaitGroup
	wg.Add(threadNo)

	go monitor.worker('A', &wg)
	go monitor.worker('B', &wg)
	go monitor.worker('C', &wg)

	wg.Wait()
	monitor.Print()
}
