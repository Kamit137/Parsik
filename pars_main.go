package main
//sdsdsdsd
import (

    "io"

    "time"
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

type Price struct {
    Price string
    Date  string
}

type Product struct {
    Name  string
    Price []Price
    Img   string
}

type User struct {
    ID       int
    Name     string
    Password string
    Email    string
    Products string // JSON строка с продуктами
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
    if r.Method == "GET" {
        url := r.URL.Query().Get("url")
        if url == "" {
            json.NewEncoder(w).Encode(map[string]string{"error": "URL is required"})
            return
        }
        product := ParseProduct(url)
        json.NewEncoder(w).Encode(product)
        return
    }
    var request struct {
        URL   string `json:"url"`
        Email string `json:"email"`
    }

    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request"})
        return
    }

    if request.URL == "" {
        json.NewEncoder(w).Encode(map[string]string{"error": "URL is required"})
        return
    }
    var product Ozon.Product
    if strings.Contains(request.URL, "ozon.ru") {
        product = Ozon.Ozon(request.URL)
    } else if strings.Contains(request.URL, "wildberries.ru") {
        product = wb_pars.Wb(request.URL)
    }
    if request.Email != "" {
        var productsJSON string
        err := db.QueryRow("SELECT products FROM users WHERE email = ?", request.Email).Scan(&productsJSON)
        if err != nil {
            json.NewEncoder(w).Encode(map[string]string{"error": "User not found"})
            return
        }

        var products []Product
        if productsJSON != "" {
            if err := json.Unmarshal([]byte(productsJSON), &products); err != nil {
                json.NewEncoder(w).Encode(map[string]string{"error": "Failed to parse products"})
                return
            }
        }


        updatedProducts, err := json.Marshal(products)
        if err != nil {
            json.NewEncoder(w).Encode(map[string]string{"error": "Failed to encode products"})
            return
        }

        _, err = db.Exec("UPDATE users SET products = ? WHERE email = ?", updatedProducts, request.Email)
        if err != nil {
            json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update database"})
            return
        }
    }

    json.NewEncoder(w).Encode(product)
}

func imageProxyHandler(w http.ResponseWriter, r *http.Request) {
    imageURL := r.URL.Query().Get("url")
    if imageURL == "" {
        http.Error(w, "URL is required", http.StatusBadRequest)
        return
    }

    client := &http.Client{Timeout: 5 * time.Second}
    resp, err := client.Get(imageURL)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
    io.Copy(w, resp.Body)
}
func ParseProduct(url string) Product {
    return Product{
        Name:  "Пример товара",
        Price: []Price{{Price: "1000", Date: time.Now().Format("2006-01-02")}},
        Img:   "https://via.placeholder.com/200",
    }
}


func registerHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    var user struct {
        Username string `json:"username"`
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }
    hashedPassword := user.Password

    _, err := db.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)",
        user.Username, user.Email, hashedPassword)

    if err != nil {
        http.Error(w, "User already exists", http.StatusConflict)
        return
    }

    w.WriteHeader(http.StatusCreated)
}

func loginHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    var credentials struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    var user User
    err := db.QueryRow("SELECT id, username, email, password, products FROM users WHERE email = ?",
        credentials.Email).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Products)

    if err != nil {
        http.Error(w, "User not found", http.StatusUnauthorized)
        return
    }

    if credentials.Password != user.Password {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }

    json.NewEncoder(w).Encode(user)
}

func main() {
    db, err := sql.Open("sqlite3", "db.db")
    if err != nil {
        fmt.Println(err)
        return
    }
    defer db.Close()

    sqlStmt := `
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            username TEXT UNIQUE NOT NULL,
            email TEXT UNIQUE NOT NULL,
            password TEXT NOT NULL,
            products TEXT DEFAULT '[]'
        );
    `
    _, err = db.Exec(sqlStmt)
    if err != nil {
        log.Printf("%q: %s\n", err, sqlStmt)
        return
    }

    http.HandleFunc("/", home)
    http.HandleFunc("/parse", enableCORS(func(w http.ResponseWriter, r *http.Request) {
        parseHandler(w, r, db)
    }))
    http.HandleFunc("/image-proxy", imageProxyHandler)
    http.HandleFunc("/register", enableCORS(func(w http.ResponseWriter, r *http.Request) {
        registerHandler(w, r, db)
    }))
    http.HandleFunc("/login", enableCORS(func(w http.ResponseWriter, r *http.Request) {
        loginHandler(w, r, db)
    }))

    http.ListenAndServe(":8092", nil)
}
