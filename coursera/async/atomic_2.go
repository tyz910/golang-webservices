package main

import (
	"fmt"
	"sync/atomic" // atomic_2.go
	"time"
)

var totalOperations int32

func inc() {
	atomic.AddInt32(&totalOperations, 1) // автомарно
}

func main() {
	for i := 0; i < 1000; i++ {
		go inc()
	}
	time.Sleep(2 * time.Millisecond)
	fmt.Println("total operation = ", totalOperations)
}
