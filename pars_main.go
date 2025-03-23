package main

import (
	"fmt"
	"net/http"

	"kod/Parsik/Site"
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Print(w, "Gooo")
}

func registration(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "registr")
}

func main() {
	a := (pars_wb.Wb("https://www.wildberries.ru/catalog/73458197/detail.aspx"))
	fmt.Print(a.Name, "\n", a.Price, "\n")
	fmt.Println(pars_wb.ScrapeImageURLs("https://www.wildberries.ru/catalog/73458197/detail.aspx"))
	http.HandleFunc("/", home)
	http.HandleFunc("/reg/", registration)
	http.ListenAndServe(":8080", nil)
}
