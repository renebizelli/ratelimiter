package middlewares_ratelimiter

import (
	"net/http"
)

type CoreInterface interface {
	Limiter(key Key, parameters *Parameters) int
	IsBlocked(blockedKey string) bool
}

type BasedonInterface interface {
	Limiter(r *http.Request, ch chan<- Response)
}
