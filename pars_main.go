package main

import (
	"fmt"
	"net/http"

	"kod/Parsik/Parser"
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Print(w, "Gooo")
}

func registration(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "registr")
}

type Biba struct{}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/reg/", registration)
	http.ListenAndServe(":8080", nil)
	parser.Pars()
}
