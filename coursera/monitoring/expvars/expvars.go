package main

import (
	"fmt"
	"net/http"
	"runtime"

	"expvar"
)

var (
	hits = expvar.NewMap("hits")
)

func handler(w http.ResponseWriter, r *http.Request) {
	hits.Add(r.URL.Path, 1)
	w.Write([]byte("expvar increased"))
}

func init() {
	expvar.Publish("mystat", expvar.Func(func() interface{} {
		hits.Init()
		return map[string]int{
			"test":          100500,
			"value":         42,
			"goroutine_num": runtime.NumGoroutine(),
		}
	}))
}

func main() {
	http.HandleFunc("/", handler)

	fmt.Println("starting server at :8081")
	http.ListenAndServe(":8081", nil)
}
