package main

import (
	"fmt"
	"database/sql"
	"log"
	_ "github.com/mattn/go-sqlite3"
	"encoding/json"
	"net/http"
	"text/template"
	"strings"
	Ozon "kod/Parsik/Sites/Ozon"
	"kod/Parsik/Sites/Wb"
)

type Price struct{
	Price string
	Date string
}
type Product struct {
	Name        string
	Price       []Price
	Img string
}


type User struct {
	ID       int
	Name     string
	Password sring
	Email    string
}

func home(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("pars4.html")
	if err != nil {
		http.Error(w, "Ошибка шаблона", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, "d")
}

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		next(w, r)
	}
}

func parseHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var url string
	if r.Method == "POST" {
		var request struct {
			URL string `json:"url"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request"})
			return
		}
		url = request.URL
	} else {
		url = r.URL.Query().Get("url")
	}

	if url == "" {
		json.NewEncoder(w).Encode(map[string]string{"error": "URL is required"})
		return
	}

var product Ozon.Product

 if strings.Contains(url, "ozon.ru") {
    product = Ozon.Ozon(url)
  } else if strings.Contains(url, "wildberries.ru") {
    product = wb_pars.Wb(url)
  }

fmt.Println(product)

	json.NewEncoder(w).Encode(product)
}

func main() {
	db, err := sql.Open("sqlite3", "db.db")
	if err != nu{
		fmt.Println(err)
		return
	}
	defer db.Close()
sqlStmt := `
		CREATE TABLE IF NOT EXISTS base (
			id INTEGER PRIMARY KEY,
			name TEXT,
			password TEXT,
			email TEXT,
			produkts TEXT
		);
	`
	_, err = db.Exec(sqlStmt) // Выполнение SQL-запроса
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt) // Выводим ошибку и сам запрос
		return
	}

	http.HandleFunc("/", home)
	http.HandleFunc("/parse", enableCORS(parseHandler))
	http.ListenAndServe(":8080", nil)
}