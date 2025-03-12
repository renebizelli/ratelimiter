package webserver

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth"
)

type LoginInput struct {
	Email                      string `json:"email"`
	RatelimiterMaxRequests     int    `json:"ratelimiterMaxRequests"`
	RaterlimiterSecondsBlocked int    `json:"raterlimiterSecondsBlocked"`
}

type LoginWebServer struct {
	JWT                               *jwtauth.JWTAuth
	JWTExpires                        int
	ratelimiterDefaultMaxRequests     int
	raterlimiterDefaultSecondsBlocked int
}

func NewLoginWebServer(jwt *jwtauth.JWTAuth, expiresIn int, ratelimiterDefaultMaxRequests int, raterlimiterDefaultSecondsBlocked int) *LoginWebServer {
	return &LoginWebServer{
		JWT:                               jwt,
		JWTExpires:                        expiresIn,
		ratelimiterDefaultMaxRequests:     ratelimiterDefaultMaxRequests,
		raterlimiterDefaultSecondsBlocked: raterlimiterDefaultSecondsBlocked,
	}
}

func (l *LoginWebServer) LoginHandler(w http.ResponseWriter, r *http.Request) {

	var input LoginInput

	e := json.NewDecoder(r.Body).Decode(&input)
	if e != nil || input.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if input.RatelimiterMaxRequests == 0 {
		input.RatelimiterMaxRequests = l.ratelimiterDefaultMaxRequests
	}

	if input.RaterlimiterSecondsBlocked == 0 {
		input.RaterlimiterSecondsBlocked = l.raterlimiterDefaultSecondsBlocked
	}

	claims := map[string]interface{}{
		"user":               input.Email,
		"rl-max-requests":    input.RatelimiterMaxRequests,
		"rl-seconds-blocked": input.RaterlimiterSecondsBlocked,
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
