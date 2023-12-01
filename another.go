package main

import (
	"fmt"
	"html/template"
	"net/http"
)

const port = ":8080"

func Home(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "Hello World")
}

func renderTemplate(w http.ResponseWriter, tmpl string) {
	t, err := template.ParseFiles("./templates/" + tmpl + ".page.html")
	if err != nil {
		fmt.Println("Error parsing template:", err)
	}
	t.Execute(w, nil)
}

func main() {
	http.HandleFunc("/", Home)

	fmt.Println("(https://localhost:8080) - Server started on port", port)
	http.ListenAndServe(":8080", nil)

	game := objects.HangmanGame{}
	game.LancerPendu()
}
