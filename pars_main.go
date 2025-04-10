package main

import (
	"fmt"
	"net/http"
	"text/template"

	"kod/Parsik/Sites/Ali"
	"kod/Parsik/Sites/Wb"
)

func home(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("less1/index2.html")
	if err != nil {
		http.Error(w, "Ошибка шаблона", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, "d")
}

func main() {
	ali.M()
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	a := wb_pars.Wb("https://www.wildberries.ru/catalog/73458197/detail.aspx")
	fmt.Println(a.Name, "\n", a.Price)
	fmt.Println(wb_pars.ScrapeImageURLs("https://www.wildberries.ru/catalog/73458197/detail.aspx"))
	http.HandleFunc("/", home)
	http.ListenAndServe(":8088", nil)
}
