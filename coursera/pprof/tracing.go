package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"runtime"
	"encoding/json"
	_ "net/http/pprof"
	"time"
)

type Post struct {
	ID       int
	Text     string
	Author   string
	Comments int
	Time     time.Time
}

func handle(w http.ResponseWriter, req *http.Request) {
	result := ""
	for i := 0; i < 100; i++ {
		currPost := &Post{ID: i, Text: "new post", Time: time.Now()}
		jsonRaw, _ := json.Marshal(currPost)
		result += string(jsonRaw)
	}
	time.Sleep(3 * time.Millisecond)
	w.Write([]byte(result))
}

func main() {
	runtime.GOMAXPROCS(4)
	http.HandleFunc("/", handle)

	fmt.Println("starting server at :8080")
	fmt.Println(http.ListenAndServe(":8080", nil))
}

/*





go build -o tracing.exe tracing.go && ./tracing.exe

ab -t 300 -n 10000000 -c 10 http://127.0.0.1:8080/

curl http://localhost:8080/debug/pprof/trace?seconds=10 -o trace.out

go tool trace -http "0.0.0.0:8081" tracing.exe trace.out


*/
