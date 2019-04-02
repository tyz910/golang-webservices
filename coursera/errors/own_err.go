package main

import (
	"fmt"
	"net/http"
	"time"
)

var (
	client = http.Client{Timeout: time.Duration(time.Millisecond)}
)

type ResourceError struct {
	URL string
	Err error
}

func (re *ResourceError) Error() string {
	return fmt.Sprintf(
		"Resource error: URL: %s, err: %v",
		re.URL,
		re.Err,
	)
}

func getRemoteResource() error {
	url := "http://127.0.0.1:9999/pages?id=123"
	_, err := client.Get(url)
	if err != nil {
		return &ResourceError{URL: url, Err: err}
	}
	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	err := getRemoteResource()
	if err != nil {
		switch err.(type) {
		case *ResourceError:
			err := err.(*ResourceError)
			fmt.Printf("resource %s err: %s\n", err.URL, err.Err)
			http.Error(w, "remote resource error", 500)
		default:
			fmt.Printf("internal error: %+v\n", err)
			http.Error(w, "internal error", 500)
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
