package webserver

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth"
)

type LoginInput struct {
	Email          string `json:"email"`
	MaxRequests    int    `json:"maxRequests"`
	BlockedSeconds int    `json:"blockedSeconds"`
}

type LoginWebServer struct {
	mux        *http.ServeMux
	JWT        *jwtauth.JWTAuth
	JWTExpires int
}

func NewLoginWebServer(mux *http.ServeMux, jwt *jwtauth.JWTAuth, expiresIn int) *LoginWebServer {

	return &LoginWebServer{
		mux:        mux,
		JWT:        jwt,
		JWTExpires: expiresIn,
	}
}

func (l *LoginWebServer) RegisterRoutes() {
	l.mux.HandleFunc("POST /login", l.loginHandler)
}

func (l *LoginWebServer) loginHandler(w http.ResponseWriter, r *http.Request) {

	var input LoginInput

	e := json.NewDecoder(r.Body).Decode(&input)
	if e != nil || input.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	claims := map[string]interface{}{
		"user":               input.Email,
		"rl-max-requests":    input.MaxRequests,
		"rl-seconds-blocked": input.BlockedSeconds,
		"exp":                time.Now().Add(time.Minute * time.Duration(l.JWTExpires)).Unix(),
	}

	log.Println(claims)

	_, stringToken, jwterr := l.JWT.Encode(claims)

	if jwterr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	json.NewEncoder(w).Encode(stringToken)
}
