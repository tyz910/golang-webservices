package main

import (
	"bytes"
	"flag"
	"fmt"
	"net"
	"net/http"
	"runtime"
	"time"

	_ "expvar"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
}

var carbonAddr = flag.String("graphite", "192.168.99.100:2003", "The address of carbon receiver")

func main() {
	flag.Parse()

	go sendStat()

	http.HandleFunc("/", handler)

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}

func sendStat() {
	m := &runtime.MemStats{}
	conn, err := net.Dial("tcp", *carbonAddr)
	if err != nil {
		panic(err)
	}
	c := time.Tick(time.Minute)
	for tickTime := range c {
		runtime.ReadMemStats(m)

		buf := bytes.NewBuffer([]byte{})
		fmt.Fprintf(buf, "coursera.mem_heap %d %d\n",
			m.HeapInuse, tickTime.Unix())
		fmt.Fprintf(buf, "coursera.mem_stack %d %d\n",
			m.StackInuse, tickTime.Unix())
		fmt.Fprintf(buf, "coursera.goroutines_num %d %d\n",
			runtime.NumGoroutine(), tickTime.Unix())

		conn.Write(buf.Bytes())
		fmt.Println(buf.String())
	}
}

/*
docker run -d\
 --name graphite\
 --restart=always\
 -p 80:80\
 -p 2003-2004:2003-2004\
 -p 2023-2024:2023-2024\
 -p 8125:8125/udp\
 -p 8126:8126\
 graphiteapp/graphite-statsd
*/
