package main

import (
	"fmt"
	"net/http"
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
	s := ""
	for i := 0; i < 1000; i++ {
		p := &Post{ID: i, Text: "new post"}
		s += fmt.Sprintf("%#v", p)
	}
	w.Write([]byte(s))
}

func main() {
	http.HandleFunc("/", handle)

	fmt.Println("starting server at :8080")
	fmt.Println(http.ListenAndServe(":8080", nil))
}

/*
go build -o pprof_1.exe pprof_1.go && ./pprof_1.exe

ab -t 300 -n 1000000000 -c 10 http://127.0.0.1:8080/

curl http://127.0.0.1:8080/debug/pprof/heap -o mem_out.txt
curl http://127.0.0.1:8080/debug/pprof/profile?seconds=5 -o cpu_out.txt

go tool pprof -svg -inuse_space pprof_1.exe mem_out.txt > mem_is.svg
go tool pprof -svg -inuse_objects pprof_1.exe mem_out.txt > mem_oo.svg
go tool pprof -svg -alloc_space pprof_1.exe mem_out.txt > mem_as.svg
go tool pprof -svg -alloc_objects pprof_1.exe mem_out.txt > mem_ao.svg
go tool pprof -svg pprof_1.exe cpu_out.txt > cpu.svg

*/
