package main

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"net/http"
)

var db *sql.DB

var err error

func main() {
	db, err = sql.Open("pgx", "postgres://app:app@localhost:5433/appdb?sslmode=disable")
	if err := db.Ping(); err != nil {
		panic(err)
	}

	fmt.Println("server is starting...")
	fmt.Println("listening on localhost:8080")
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register", registerHandler)
	err = http.ListenAndServe(":8080", nil)
	panic(err)

}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Alles gut")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method is not allowed", http.StatusMethodNotAllowed)
		return
	}
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")
	if username == "" || password == "" {
		http.Error(w, "wrong username or password", http.StatusBadRequest)
	}
	var dbPassword string
	err := db.QueryRow("SELECT password FROM users WHERE username = $1", username).Scan(&dbPassword)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	if dbPassword != password {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	fmt.Fprintln(w, "Login is ok.")
}
func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method is not allowed", http.StatusMethodNotAllowed)
		return
	}
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")
	if username == "" || password == "" {
		http.Error(w, "username or password is empty", http.StatusBadRequest)
		return
	}
	_, err := db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", username, password)
	if err != nil {
		http.Error(w, "db error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "User registered successfully.")
}
