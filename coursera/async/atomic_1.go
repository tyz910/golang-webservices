package main

import (
	"fmt"
	"time"
)

var totalOperations int32 = 0

func inc() {
	totalOperations++
}

func main() {
	for i := 0; i < 1000; i++ {
		go inc()
	}
	time.Sleep(2 * time.Millisecond)
	// ождается 1000, но по факту будет меньше
	fmt.Println("total operation = ", totalOperations)
}
