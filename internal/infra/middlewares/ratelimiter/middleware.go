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

		var key Key
		var parameters *Parameters

		if l.tokenOn {
			err := l.basedOnToken.Validate(r)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"message" : "` + err.Error() + `"}`))
				return
			}

			key, parameters = l.basedOnToken.Parse(r)

		} else if l.ipOn {
			err := l.basedOnIP.Validate()
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"message" : "` + err.Error() + `"}`))
				return
			}

			key, parameters = l.basedOnIP.Parse(r)
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

		next.ServeHTTP(w, r)
	})
}
