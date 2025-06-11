package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/rodrigocitadin/insecure/models"
)

var db *sql.DB

func InitDB(database *sql.DB) {
	db = database
	db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		password TEXT,
		amount REAL DEFAULT 100.0
	);`)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var u models.User
	json.NewDecoder(r.Body).Decode(&u)
	db.Exec("INSERT INTO users (name, password) VALUES (?, ?)", u.Name, u.Password)
	w.WriteHeader(201)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var u models.User
	json.NewDecoder(r.Body).Decode(&u)
	row := db.QueryRow("SELECT id FROM users WHERE name = ? AND password = ?", u.Name, u.Password)
	var id int
	if err := row.Scan(&id); err != nil {
		http.Error(w, "invalid credentials", 401)
		return
	}
	w.Write([]byte("user-" + strconv.Itoa(id)))
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	row := db.QueryRow("SELECT id, name, amount FROM users WHERE id = ?", id)
	var u models.User
	if err := row.Scan(&u.ID, &u.Name, &u.Balance); err != nil {
		http.Error(w, "user not found", 404)
		return
	}
	json.NewEncoder(w).Encode(u)
}

func DumpHandler(w http.ResponseWriter, r *http.Request) {
	rows, _ := db.Query("SELECT id, name, pass, amount FROM users")
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		rows.Scan(&u.ID, &u.Name, &u.Password, &u.Balance)
		users = append(users, u)
	}
	json.NewEncoder(w).Encode(users)
}

func BalanceHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Authorization")
	fmt.Println(userID)
	row := db.QueryRow("SELECT amount FROM users WHERE id = ?", userID[len("user-"):])
	var amount float64
	if err := row.Scan(&amount); err != nil {
		http.Error(w, "user not found", 404)
		return
	}
	w.Write([]byte("Saldo: " + formatFloat(amount)))
}
