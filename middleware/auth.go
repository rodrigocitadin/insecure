package middleware

import (
	"net/http"
	"sync"
)

var rateLimit = make(map[string]int)
var rateMutex sync.Mutex

func WithCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next(w, r)
	}
}

func RateLimiter(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		rateMutex.Lock()
		rateLimit[ip]++
		count := rateLimit[ip]
		rateMutex.Unlock()

		if count > 10 {
			http.Error(w, "Too many requests", 429)
			return
		}
		next(w, r)
	}
}
