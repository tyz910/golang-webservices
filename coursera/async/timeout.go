package main

import (
	"fmt"
	"time"
)

func longSQLQuery() chan bool {
	ch := make(chan bool, 1)
	go func() {
		time.Sleep(2 * time.Second)
		ch <- true
	}()
	return ch
}

func main() {
	// при 1 выполнится таймаут, при 3 - выполнится операция
	timer := time.NewTimer(3 * time.Second)
	select {
	case <-timer.C:
		fmt.Println("timer.C timeout happend")
	case <-time.After(time.Minute):
		// пока не выстрелит - не соберётся сборщиком мусора
		fmt.Println("time.After timeout happend")
	case result := <-longSQLQuery():
		// освобождет ресурс
		if !timer.Stop() {
			<-timer.C
		}
		fmt.Println("operation result:", result)
	}
}
