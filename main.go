package main

import (
	// "fmt"
	"html/template"
	"net/http"
)

var deja = []string{}

type Todo struct {
	Title string
	Done  bool
}

type Hangman struct {
	Deja []string
}

func main() {
	tmpl := template.Must(template.ParseFiles("layout.html"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := Hangman{
			Deja: deja,
		}

		lettre := r.FormValue("lettre")
		if lettre != "" {
			deja = append(deja, lettre)
		}

		tmpl.Execute(w, data)
	})

	http.ListenAndServe(":80", nil)
}
