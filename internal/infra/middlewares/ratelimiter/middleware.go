package middlewares_ratelimiter

import (
	"net/http"
)

type RateLimiterMiddleware struct {
	basedOnToken *BasedOnToken
	basedOnIP    *BasedOnIP
	core         CoreInterface
	ipOn         bool
	tokenOn      bool
}

func NewRateLimiterMiddleware(
	ipOn bool,
	tokenOn bool,
	basedOnToken *BasedOnToken,
	basedOnIP *BasedOnIP,
	core CoreInterface,
) *RateLimiterMiddleware {

	return &RateLimiterMiddleware{
		ipOn:         ipOn,
		tokenOn:      tokenOn,
		basedOnToken: basedOnToken,
		basedOnIP:    basedOnIP,
		core:         core,
	}
}

var message409 = "you have reached the maximum number of requests or actions allowed within a certain time frame"

func (l *RateLimiterMiddleware) Limiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if l.tokenOn || l.ipOn {

			var key Key
			var parameters *Parameters
			var e error

			if l.tokenOn {
				key, parameters, e = l.basedOnToken.Parse(r)
			} else if l.ipOn {
				key, parameters, e = l.basedOnIP.Parse(r)
			}

			if e != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"message" : "` + e.Error() + `"}`))
				return
			}

			httpStatus := l.core.Limiter(r.Context(), key, parameters)

			if httpStatus == http.StatusTooManyRequests {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte(`{"message" : "` + message409 + `"}`))
				return
			} else if httpStatus != http.StatusOK {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(int(httpStatus))
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
