package webserver

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth"
)

type LoginInput struct {
	Email string `json:"email"`
}

type LoginWebServer struct {
	JWT        *jwtauth.JWTAuth
	JWTExpires int64
}

func NewLoginWebServer(jwt *jwtauth.JWTAuth, expiresIn int64) *LoginWebServer {
	return &LoginWebServer{
		JWT:        jwt,
		JWTExpires: expiresIn,
	}
}

func (l *LoginWebServer) LoginHandler(w http.ResponseWriter, r *http.Request) {

	var input LoginInput

	e := json.NewDecoder(r.Body).Decode(&input)
	if e != nil || input.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	claims := map[string]interface{}{
		"sub": input.Email,
		"exp": time.Now().Add(time.Second * time.Duration(l.JWTExpires)).Unix(),
	}

	_, stringToken, jwterr := l.JWT.Encode(claims)

	if jwterr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	json.NewEncoder(w).Encode(stringToken)
}
