package main

import (
	"fmt"
	"net/http"
	"time"
)

var (
	client = http.Client{Timeout: time.Duration(time.Millisecond)}
)

func getRemoteResource() error {
	url := "http://127.0.0.1:9999/pages?id=123"
	_, err := client.Get(url)
	if err != nil {
		// вернётся `timed out`. и что?
		// return err

		// будет `res error: time out`. а где?
		// return fmt.Errof("getRemoteResource: %+v", err)

		return fmt.Errorf("getRemoteResource: %s at %s", err, url)
	}
	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	err := getRemoteResource()

	if err != nil {
		fmt.Printf("error happend: %+v\n", err)
		http.Error(w, "internal error", 500)
		return
	}

	w.Write([]byte("all is OK"))
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}
