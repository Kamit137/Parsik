package main

import (
	"encoding/json"
	// "fmt"
	"net/http"
	"text/template"
	// ss
	// "kod/Parsik/Sites/Ali"
	"kod/Parsik/Sites/Wb"
)

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

	product := wb_pars.Wb(url)

	response := map[string]interface{}{
		"name":     product.Name,
		"price":    product.Price,
		"imageUrl": product.ImageURL,
		"rating":   4.5,
		"seller":   "Wildberries",
	}

	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/parse", enableCORS(parseHandler))
	http.ListenAndServe(":8080", nil)
}
