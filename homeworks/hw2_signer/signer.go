package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

// ExecutePipeline обеспечивает конвейерную обработку функций-воркеров
func ExecutePipeline(jobs ...job) {
	var in, out chan interface{}
	wg := &sync.WaitGroup{}

	for _, j := range jobs {
		in = out
		out = make(chan interface{}, 100)

		wg.Add(1)
		go func(j job, in, out chan interface{}) {
			defer wg.Done()
			defer close(out)

			j(in, out)
		}(j, in, out)
	}

	wg.Wait()
}

// processStr Параллельно читает строки из in, обрабатывает через функцию process и записывает в out
func processStr(in, out chan interface{}, process func(string) string) {
	wg := &sync.WaitGroup{}

	for data := range in {
		wg.Add(1)

		go func(data string) {
			defer wg.Done()

			out <- process(data)
		}(fmt.Sprintf("%v", data))
	}

	wg.Wait()
}

// SingleHash считает значение crc32(data)+"~"+crc32(md5(data)) ( конкатенация двух строк через ~),
// где data - то что пришло на вход (по сути - числа из первой функции)
func SingleHash(in, out chan interface{}) {
	md5Mux := &sync.Mutex{}

	processStr(in, out, func(data string) string {
		// crc32(data)
		hash1 := make(chan string)
		go func() {
			defer close(hash1)

			hash1 <- DataSignerCrc32(data)
		}()

		// crc32(md5(data))
		hash2 := make(chan string)
		go func() {
			defer close(hash2)

			// DataSignerMd5 может одновременно вызываться только 1 раз
			md5Mux.Lock()
			md5Hash := DataSignerMd5(data)
			md5Mux.Unlock()

			hash2 <- DataSignerCrc32(md5Hash)
		}()

		return fmt.Sprintf("%s~%s", <-hash1, <-hash2)
	})
}

// MultiHash считает значение crc32(th+data) (конкатенация цифры, приведённой к строке и строки),
// где th=0..5 (т.е. 6 хешей на каждое входящее значение), потом берёт конкатенацию результатов в порядке расчета (0..5),
// где data - то что пришло на вход (и ушло на выход из SingleHash)
func MultiHash(in, out chan interface{}) {
	processStr(in, out, func(data string) string {
		results := make([]string, 6)
		resMux := &sync.Mutex{}
		resWg := &sync.WaitGroup{}

		for th := 0; th <= 5; th++ {
			resWg.Add(1)

			go func(th int) {
				defer resWg.Done()

				hash := DataSignerCrc32(fmt.Sprintf("%d%s", th, data)) // crc32(th+data)

				// Конкурентная запись
				resMux.Lock()
				results[th] = hash
				resMux.Unlock()
			}(th)
		}

		resWg.Wait()
		return strings.Join(results, "")
	})
}

// CombineResults получает все результаты, сортирует (https://golang.org/pkg/sort/),
// объединяет отсортированный результат через _ (символ подчеркивания) в одну строку
func CombineResults(in, out chan interface{}) {
	var results []string

	for data := range in {
		results = append(results, fmt.Sprintf("%v", data))
	}

	sort.Strings(results)
	out <- strings.Join(results, "_")
}
