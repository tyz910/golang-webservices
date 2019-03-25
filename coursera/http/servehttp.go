package main

import (
	"fmt"
	"net/http"
)

type Handler struct {
	Name string
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Name:", h.Name, "URL:", r.URL.String())
}

func main() {
	testHandler := &Handler{Name: "test"}
	http.Handle("/test/", testHandler)

	rootHandler := &Handler{Name: "root"}
	http.Handle("/", rootHandler)

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}
