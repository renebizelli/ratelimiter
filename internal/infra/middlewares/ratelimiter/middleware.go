package middlewares_ratelimiter

import (
	"net/http"
)

type RateLimiterMiddleware struct {
	basedOns []BasedonInterface
}

func NewRateLimiterMiddleware(
	basedOns []BasedonInterface,
) *RateLimiterMiddleware {

	return &RateLimiterMiddleware{
		basedOns: basedOns,
	}
}

var message409 = "you have reached the maximum number of requests or actions allowed within a certain time frame"

func (l *RateLimiterMiddleware) Limiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		for _, basedOn := range l.basedOns {

			e := basedOn.Limiter(r)

			if e != nil {

				if e.HttpStatus == http.StatusTooManyRequests {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusTooManyRequests)
					w.Write([]byte(`{"message" : "` + message409 + `"}`))
					return
				} else if e.HttpStatus != http.StatusOK {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(e.HttpStatus)
					w.Write([]byte(`{"message" : "` + e.Error() + `"}`))
					return
				}

			}
		}

		next.ServeHTTP(w, r)
	})
}
