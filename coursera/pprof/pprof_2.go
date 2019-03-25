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

func getPost(out chan []Post) {
	posts := []Post{}
	for i := 1; i < 10; i++ {
		post := Post{ID: 1, Text: "text"}
		posts = append(posts, post)
	}
	out <- posts
}

func handleLeak(w http.ResponseWriter, req *http.Request) {
	res := make(chan []Post)
	go getPost(res)
}

func main() {
	http.HandleFunc("/", handleLeak)

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}

/*

go build -o pprof_2.exe pprof_2.go && ./pprof_2.exe

ab -n 1000 -c 10 http://127.0.0.1:8080/

curl http://localhost:8080/debug/pprof/goroutine?debug=2 -o goroutines.txt

*/
