package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"gopkg.in/alexcesaro/statsd.v2"
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

type TimingMiddleware struct {
	sync.Mutex
	StatsReciever *statsd.Client
	Metrics       map[string]int
}

func NewTimingMiddleware(st *statsd.Client) *TimingMiddleware {
	tm := &TimingMiddleware{
		StatsReciever: st,
		Metrics:       make(map[string]int),
	}
	return tm
}

func (tm *TimingMiddleware) TrackRequestTimings(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx,
			timingsKey,
			&ctxTimings{
				Data: make(map[string]*Timing),
			})
		defer tm.logContextTimings(ctx, r.URL.Path, time.Now())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (tm *TimingMiddleware) logContextTimings(ctx context.Context, path string, start time.Time) {
	// получаем тайминги из контекста
	// поскольку там пустой интерфейс, то нам надо преобразовать к нужному типу
	timings, ok := ctx.Value(timingsKey).(*ctxTimings)
	if !ok {
		return
	}
	totalReal := time.Since(start)

	path = strings.Replace(path, "/", "-", -1)

	prefix := "requests." + path + "."

	// buf := bytes.NewBufferString(path)
	var total time.Duration
	for timing, value := range timings.Data {
		metric := prefix + "timings." + timing
		tm.StatsReciever.Increment(metric)
		tm.StatsReciever.Timing(metric+"_time", uint64(value.Duration/time.Millisecond))
		total += value.Duration
	}

	tm.StatsReciever.Increment(prefix + "hits")
	tm.StatsReciever.Timing(prefix+"tracked", uint64(totalReal/time.Millisecond))
	tm.StatsReciever.Timing(prefix+"real_time", uint64(total/time.Millisecond))
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

var (
	statsdAddr = flag.String("addr", "192.168.99.100:32770", "statsd addr")
)

func main() {
	flag.Parse()

	rand.Seed(time.Now().UTC().UnixNano())

	siteMux := http.NewServeMux()
	siteMux.HandleFunc("/", loadPostsHandle)

	statsd, err := statsd.New(
		statsd.Prefix("prod"),
		statsd.Address(*statsdAddr),
		statsd.ErrorHandler(func(err error) {
			log.Print(err)
		}),
	)
	if err != nil {
		log.Print(err)
	}

	tm := NewTimingMiddleware(statsd)
	siteHandler := tm.TrackRequestTimings(siteMux)

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", siteHandler)
}
