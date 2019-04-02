package main

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
)

var (
	client = http.Client{Timeout: time.Duration(time.Millisecond)}
)

func getRemoteResource() error {
	url := "http://127.0.0.1:9999/pages?id=123"
	_, err := client.Get(url)
	if err != nil {
		return errors.Wrap(err, "resource error")
	}
	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	err := getRemoteResource()
	if err != nil {
		fmt.Printf("full err: %+v\n", err)
		switch err := errors.Cause(err).(type) {
		case *url.Error:
			fmt.Printf("resource %s err: %+v\n", err.URL, err.Err)
			http.Error(w, "remote resource error", 500)
		default:
			fmt.Printf("%+v\n", err)
			http.Error(w, "parsing error", 500)
		}
		return
	}
	w.Write([]byte("all is OK"))
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}
