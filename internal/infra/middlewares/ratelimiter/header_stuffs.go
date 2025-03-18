package middlewares_ratelimiter

import (
	"errors"
	"net/http"
	"strings"
)

type HeaderByStuffs struct {
}

func (h *HeaderByStuffs) Set(r *http.Request) {
	r.Header.Add("rl", "ok")
	r.Header.Add("rl-basedon", "token")
}

func (h *HeaderByStuffs) IsByPass(r *http.Request) bool {

	if v := r.Header.Get("rl"); v == "ok" {
		return true
	}

	return false
}

func (h *HeaderByStuffs) GetAPIKey(r *http.Request) (string, error) {

	token := strings.Trim(r.Header.Get("API_KEY"), " ")

	if token == "" {
		return "", errors.New("invalid API_KEY")
	}

	return token, nil
}

func (h *HeaderByStuffs) SetAPIKey(r *http.Request, value string) {
	r.Header.Add("API_KEY", value)
}
