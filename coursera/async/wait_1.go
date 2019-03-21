package main

import (
	"log"
	"time"
)

func main() {
	result := make(chan string)
	go func(out chan<- string) {
		time.Sleep(1 * time.Second)
		log.Println("async operation ready, return result")
		out <- "success"
	}(result)

	time.Sleep(2 * time.Second)
	log.Println("some userful work")

	opStatus := <-result
	log.Println("main goroutine:", opStatus)
}
