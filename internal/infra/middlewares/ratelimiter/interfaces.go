package middlewares_ratelimiter

import (
	"context"
	"net/http"
)

type CoreInterface interface {
	Limiter(ctx context.Context, key Key, parameters *Parameters) int
}

type BasedonInterface interface {
	Limiter(r *http.Request, ch chan<- Response)
}
