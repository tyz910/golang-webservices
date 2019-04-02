package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

// сколько в среднем спим при эмуляции работы
const AvgSleep = 50

func trackContextTimings(ctx context.Context, metricName string, start time.Time) {
	// получаем тайминги из контекста
	// поскольку там пустой интерфейс, то нам надо преобразовать к нужному типу
	timings, ok := ctx.Value(timingsKey).(*ctxTimings)
	if !ok {
		return
	}
	elapsed := time.Since(start)
	// лочимся на случай конкурентной записи в мапку
	timings.Lock()
	defer timings.Unlock()
	// если меткри ещё нет - мы её создадим, если есть - допишем в существующую
	if metric, metricExist := timings.Data[metricName]; !metricExist {
		timings.Data[metricName] = &Timing{
			Count:    1,
			Duration: elapsed,
		}
	} else {
		metric.Count++
		metric.Duration += elapsed
	}
}

type Timing struct {
	Count    int
	Duration time.Duration
}

type ctxTimings struct {
	sync.Mutex
	Data map[string]*Timing
}

// линтер ругается если используем базовые типы в Value контекста
// типа так безопаснее разграничивать
type key int

const timingsKey key = 1

func logContextTimings(ctx context.Context, path string, start time.Time) {
	// получаем тайминги из контекста
	// поскольку там пустой интерфейс, то нам надо преобразовать к нужному типу
	timings, ok := ctx.Value(timingsKey).(*ctxTimings)
	if !ok {
		return
	}
	totalReal := time.Since(start)
	buf := bytes.NewBufferString(path)
	var total time.Duration
	for timing, value := range timings.Data {
		total += value.Duration
		buf.WriteString(fmt.Sprintf("\n\t%s(%d): %s", timing, value.Count, value.Duration))
	}
	buf.WriteString(fmt.Sprintf("\n\ttotal: %s", totalReal))
	buf.WriteString(fmt.Sprintf("\n\ttracked: %s", total))
	buf.WriteString(fmt.Sprintf("\n\tunkn: %s", totalReal-total))

	fmt.Println(buf.String())
}

func timingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx,
			timingsKey,
			&ctxTimings{
				Data: make(map[string]*Timing),
			})
		defer logContextTimings(ctx, r.URL.Path, time.Now())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func emulateWork(ctx context.Context, workName string) {
	defer trackContextTimings(ctx, workName, time.Now())

	rnd := time.Duration(rand.Intn(AvgSleep))
	time.Sleep(time.Millisecond * rnd)
}

func loadPostsHandle(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	emulateWork(ctx, "checkCache")
	emulateWork(ctx, "loadPosts")
	emulateWork(ctx, "loadPosts")
	emulateWork(ctx, "loadPosts")
	time.Sleep(10 * time.Millisecond)
	emulateWork(ctx, "loadSidebar")
	emulateWork(ctx, "loadComments")

	fmt.Fprintln(w, "Request done")
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	siteMux := http.NewServeMux()
	siteMux.HandleFunc("/", loadPostsHandle)

	siteHandler := timingMiddleware(siteMux)

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", siteHandler)
}
