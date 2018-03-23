package main

import (
	"sync"
)

var sum struct {
	sync.Mutex
	count float64
}

func CounterUp() {
	sum.Lock()
	sum.count++
	sum.Unlock()
}

func CounterReset() float64 {
	sum.Lock()
	value := sum.count
	sum.count = 0
	sum.Unlock()
	return value
}
