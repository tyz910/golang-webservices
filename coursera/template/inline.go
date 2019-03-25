package main

import (
	"fmt"
	"net/http"
	"text/template"
)

type tplParams struct {
	URL     string
	Browser string
}

const EXAMPLE = `
Browser {{.Browser}}

you at {{.URL}}
`

func handle(w http.ResponseWriter, r *http.Request) {
	tmpl := template.New(`example`)
	tmpl, _ = tmpl.Parse(EXAMPLE)

	params := tplParams{
		URL:     r.URL.String(),
		Browser: r.UserAgent(),
	}

	tmpl.Execute(w, params)
}

func main() {
	http.HandleFunc("/", handle)

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}
