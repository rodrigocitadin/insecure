package handlers

import (
	"encoding/json"
	"html/template"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/rodrigocitadin/insecure/models"
)

func TransferHandler(w http.ResponseWriter, r *http.Request) {
	sourceID := r.Header.Get("Authorization")[len("user-"):]
	targetID := r.URL.Query().Get("to")
	amountStr := r.URL.Query().Get("amount")
	amount, _ := strconv.ParseFloat(amountStr, 64)

	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)

	var current float64
	db.QueryRow("SELECT amount FROM users WHERE id = ?", sourceID).Scan(&current)

	if current < amount {
		http.Error(w, "not enough balance", 400)
		return
	}

	db.Exec("UPDATE users SET amount = amount - ? WHERE id = ?", amount, sourceID)
	db.Exec("UPDATE users SET amount = amount + ? WHERE id = ?", amount, targetID)

	w.Write([]byte("transfer made successfully"))
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	query := "SELECT id, name FROM users WHERE name LIKE '%" + name + "%'"
	rows, _ := db.Query(query)
	defer rows.Close()

	var results []models.User
	for rows.Next() {
		var u models.User
		rows.Scan(&u.ID, &u.Name)
		results = append(results, u)
	}
	json.NewEncoder(w).Encode(results)
}

func GreetHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	tmpl := `<html><body><h1>Olá, ` + name + `!</h1></body></html>`
	w.Header().Set("Content-Type", "text/html")
	template.Must(template.New("greet").Parse(tmpl)).Execute(w, nil)
}

func formatFloat(v float64) string {
	return strconv.FormatFloat(v, 'f', 2, 64)
}
