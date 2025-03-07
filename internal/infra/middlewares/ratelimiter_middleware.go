package middlewares

import (
	"log"
	"net/http"
)

func RateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Custom logger: ", r.Method, r.URL.Path)

		//w.Header().Set("Content-Type", "application/json")
		//w.WriteHeader(429)

		next.ServeHTTP(w, r)
	})
}
