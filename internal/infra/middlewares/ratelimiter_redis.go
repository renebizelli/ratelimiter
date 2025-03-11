package middlewares

import (
	"log"
	"net/http"

	"github.com/go-chi/jwtauth"
	"github.com/renebizelli/ratelimiter/configs"
)

func RateLimiter(configs *configs.Config, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenString := jwtauth.TokenFromHeader(r)

		if tokenString == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			w.Write([]byte(`{"message": "invalid authorization header"}`))
			return
		}

		_, e := jwtauth.VerifyToken(configs.JWTToken, tokenString)

		if e != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			w.Write([]byte(`{"message": "authorization token invalid"}`))
			return
		}

		token, _ := configs.JWTToken.Decode(tokenString)
		sub, _ := token.Get("sub")

		log.Println("Authorization: ", token)
		log.Println("claims: ", sub)

		log.Println("claims: ", r.RemoteAddr)

		//w.Header().Set("Content-Type", "application/json")
		//w.WriteHeader(429)

		next.ServeHTTP(w, r)
	})
}
