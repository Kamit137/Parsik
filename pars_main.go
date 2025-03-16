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

func main() {
	fmt.Println(pars_wb.Wb("https://www.wildberries.ru/catalog/187274636/detail.aspx"))
	http.HandleFunc("/", home)
	http.HandleFunc("/reg/", registration)
	http.ListenAndServe(":8080", nil)
}
