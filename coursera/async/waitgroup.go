package main

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	iterationsNum = 7
	goroutinesNum = 5
)

func startWorker(in int, wg *sync.WaitGroup) {
	defer wg.Done() // wait_2.go уменьшаем счетчик на 1
	for j := 0; j < iterationsNum; j++ {
		fmt.Printf(formatWork(in, j))
		runtime.Gosched()
	}
}

func main() {
	wg := &sync.WaitGroup{} // wait_2.go инициализируем группу
	for i := 0; i < goroutinesNum; i++ {
		wg.Add(1) // wait_2.go добавляем воркер
		go startWorker(i, wg)
	}
	time.Sleep(time.Millisecond)

	// fmt.Scanln()
	wg.Wait() // wait_2.go ожидаем, пока waiter.Done() не приведёт счетчик к 0
}

func formatWork(in, j int) string {
	return fmt.Sprintln(strings.Repeat("  ", in), "█",
		strings.Repeat("  ", goroutinesNum-in),
		"th", in,
		"iter", j, strings.Repeat("■", j))
}
