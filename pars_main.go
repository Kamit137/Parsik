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
	Ozon "Parsik/Sites/Ozon"
	"Parsik/Sites/Wb"
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
	Password string
	email    string
}

func home(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("pars4.html")
	if err != nil {
		http.Error(w, "Ошибка шаблона", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, "")
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

func parseHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")

	var url string
	var email string
	if r.Method == "POST" {
		var request struct {
			URL string `json:"url"`
			email string `json:"email"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request"})
			return
		}
		url = request.URL
		email = request.email

	} else {
		url = r.URL.Query().Get("url")
		email = r.URL.Query().Get("email")
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

var proshed_JSON string

err := db.QueryRow("SELECT produkts FROM base WHERE email = ?", email).Scan(&proshed_JSON)
	if err != nil {
		fmt.Println("ошибка", err, email)
}
var array []interface{}
	if proshed_JSON != "" {
		err = json.Unmarshal([]byte(proshed_JSON), &array)

	} else {
		array = make([]interface{}, 0) // Инициализируем пустой массив, если поле пустое
	}

	// 3. Декодирование нового JSON-объекта
	p,err := json.Marshal(product)
	if err != nil{
		fmt.Println(err,"err")
	}
	var newObject map[string]interface{}
	err = json.Unmarshal([]byte(p), &newObject)
	if err != nil {
		fmt.Printf("ошибка при декодировании нового JSON-объекта: %w", err)
	}

	// 4. Добавление нового объекта в массив
	array = append(array, newObject)

	// 5. Кодирование объединенного JSON-массива
	combinedJSON, err := json.Marshal(array)
	if err != nil {
		fmt.Printf("ошибка при кодировании объединенного JSON: %w", err)
	}

	// 6. Запись объединенного JSON-массива в базу данных
	_, err = db.Exec("UPDATE base SET produkts = ? WHERE id = ?", combinedJSON, id)
	if err != nil {
		 fmt.Printf("ошибка при записи объединенного JSON в базу данных: %w", err)
	}

	fmt.Println("JSON успешно обновлен.")


}




func main() {
	db, err := sql.Open("sqlite3", "db.db")
	if err != nil{
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
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt) // Выводим ошибку и сам запрос
		return
	}

	http.HandleFunc("/", home)
	http.HandleFunc("/parse", enableCORS(parseHandler),db)
	http.ListenAndServe(":8080", nil)
}