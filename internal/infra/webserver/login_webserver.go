package webserver

import (
	"encoding/json"
	"net/http"
	"time"

	pkg_utils "github.com/renebizelli/ratelimiter/pkg/utils"
)

type LoginInput struct {
	Email          string `json:"email"`
	MaxRequests    int    `json:"maxRequests"`
	BlockedSeconds int    `json:"blockedSeconds"`
}

type LoginWebServer struct {
	mux        *http.ServeMux
	jwt        *pkg_utils.Jwt
	jwtExpires int
}

func NewLoginWebServer(mux *http.ServeMux, jwt *pkg_utils.Jwt, expiresIn int) *LoginWebServer {

	return &LoginWebServer{
		mux:        mux,
		jwt:        jwt,
		jwtExpires: expiresIn,
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
		"key":                input.Email,
		"rl-max-requests":    input.MaxRequests,
		"rl-seconds-blocked": input.BlockedSeconds,
		"exp":                time.Now().Add(time.Minute * time.Duration(l.jwtExpires)).Unix(),
	}

	stringToken, jwterr := l.jwt.Generate(claims)

	if jwterr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	json.NewEncoder(w).Encode(stringToken)
}
