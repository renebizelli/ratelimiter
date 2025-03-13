package middlewares_ratelimiter

import "context"

type CoreInterface interface {
	Limiter(ctx context.Context, key Key, parameters *Parameters) HttpStatus
}
