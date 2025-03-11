package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/jwtauth"
	"github.com/redis/go-redis/v9"
	"github.com/renebizelli/ratelimiter/configs"
)

type RateLimiter struct {
	db     *redis.Client
	config *configs.Config
}

func NewRateLimiter(db *redis.Client, config *configs.Config) *RateLimiter {
	return &RateLimiter{
		db:     db,
		config: config,
	}
}

func (l *RateLimiter) Limiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenString := jwtauth.TokenFromHeader(r)

		if tokenString == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			w.Write([]byte(`{"message": "invalid authorization header"}`))
			return
		}

		_, e := jwtauth.VerifyToken(l.config.JWTToken, tokenString)

		if e != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			w.Write([]byte(`{"message": "authorization token invalid"}`))
			return
		}

		token, _ := l.config.JWTToken.Decode(tokenString)
		sub, _ := token.Get("sub")

		counterKey := fmt.Sprintf("%s:%d", sub, time.Now().Second())
		blockedKey := fmt.Sprintf("%s:blocked", sub)

		blocked, _ := l.db.Get(r.Context(), blockedKey).Result()

		if blocked != "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(429)
			w.Write([]byte(`{"message": "blocked"}`))
			return
		}

		count, e := l.db.Get(r.Context(), counterKey).Result()

		if e != nil {
			fmt.Println("2 >>", e)
			l.db.Set(r.Context(), counterKey, "1", time.Duration(10*time.Second))
			next.ServeHTTP(w, r)
			return
		}

		intCount, _ := strconv.Atoi(count)

		if intCount > 3 {
			log.Println("blocked")
			l.db.Set(r.Context(), blockedKey, 0, time.Duration(10*time.Second))
			w.WriteHeader(429)
			return
		}

		l.db.Set(r.Context(), counterKey, intCount+1, time.Duration(10*time.Second))

		log.Println("claims: ", counterKey)

		log.Println("RemoteAddr: ", r.RemoteAddr)
		log.Println("intValue: ", intCount)

		//w.Header().Set("Content-Type", "application/json")
		//w.WriteHeader(429)

		next.ServeHTTP(w, r)
	})
}
