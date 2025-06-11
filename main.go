package main

import (
	"database/sql"
	"log"
	"net/http"
	"sync"

	"github.com/rodrigocitadin/insecure/handlers"
	"github.com/rodrigocitadin/insecure/middleware"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var rateLimit = make(map[string]int)
var rateMutex sync.Mutex

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./insecure.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	handlers.InitDB(db)

	http.HandleFunc("/register", middleware.WithCORS(middleware.RateLimiter(handlers.RegisterHandler)))
	http.HandleFunc("/login", middleware.WithCORS(middleware.RateLimiter(handlers.LoginHandler)))
	http.HandleFunc("/balance", middleware.WithCORS(middleware.RateLimiter(handlers.BalanceHandler)))
	http.HandleFunc("/transfer", middleware.WithCORS(middleware.RateLimiter(handlers.TransferHandler)))
	http.HandleFunc("/search", middleware.WithCORS(middleware.RateLimiter(handlers.SearchHandler)))
	http.HandleFunc("/greet", middleware.WithCORS(middleware.RateLimiter(handlers.GreetHandler)))
	http.HandleFunc("/dump", middleware.WithCORS(handlers.DumpHandler))
	http.HandleFunc("/user", middleware.WithCORS(handlers.UserHandler))

	log.Println("API running on :8080")
	http.ListenAndServe(":8080", nil)
}
