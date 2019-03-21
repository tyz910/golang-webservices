package main

import (
	"fmt"
)

func main() {
	in := make(chan int, 0)

	go func(out chan<- int) {
		for i := 0; i <= 10; i++ {
			fmt.Println("before", i)
			out <- i
			fmt.Println("after", i)
		}
		close(out)
		fmt.Println("generator finish")
	}(in)

	for i := range in {
		fmt.Println("\tget", i)
	}

	// fmt.Scanln()
}
