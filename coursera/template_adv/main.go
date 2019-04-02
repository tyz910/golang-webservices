package main

import (
	"bytes"
	"fmt"
	"net/http"

	"coursera/template_adv/item"
	"coursera/template_adv/template"
)

//go:generate hero -source=./template/

var ExampleItems = []*item.Item{
	&item.Item{1, "rvasily", "Mail.ru Group"},
	&item.Item{2, "username", "freelancer"},
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		buffer := new(bytes.Buffer)
		template.Index(ExampleItems, buffer)
		w.Write(buffer.Bytes())
	})

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}
