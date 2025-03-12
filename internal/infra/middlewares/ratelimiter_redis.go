package middlewares

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/jwtauth"
	"github.com/redis/go-redis/v9"
	pkg_utils "github.com/renebizelli/ratelimiter/pkg/utils"
)

type RateLimiter struct {
	db                          *redis.Client
	jwtauth                     *jwtauth.JWTAuth
	ratelimiterIpOn             bool
	ratelimiterIpMaxRequests    int
	ratelimiterIpSecondsBlocked int
	ratelimiterTokenOn          bool
}

type CustomError struct {
	HttpStatus int
	Message    string
}

func (e *CustomError) Error() string {
	return e.Message
}

func NewRateLimiter(db *redis.Client,
	jwtauth *jwtauth.JWTAuth,
	ratelimiterIpOn bool,
	ratelimiterIpMaxRequests int,
	ratelimiterIpSecondsBlocked int,
	ratelimiterTokenOn bool,
) *RateLimiter {

	e := parametersValidate(
		ratelimiterIpMaxRequests,
		ratelimiterIpSecondsBlocked)

	pkg_utils.PanicIfError(e, "rate limiter parameters error")

	return &RateLimiter{
		db:                          db,
		jwtauth:                     jwtauth,
		ratelimiterIpOn:             ratelimiterIpOn,
		ratelimiterIpMaxRequests:    ratelimiterIpMaxRequests,
		ratelimiterIpSecondsBlocked: ratelimiterIpSecondsBlocked,
		ratelimiterTokenOn:          ratelimiterTokenOn,
	}
}

func parametersValidate(
	ratelimiterIpMaxRequests int,
	ratelimiterIpSecondsBlocked int) error {

	if ratelimiterIpMaxRequests == 0 {
		return errors.New("RATELIMITER_IP_MAX_REQUESTS is required")
	} else if ratelimiterIpSecondsBlocked == 0 {
		return errors.New("RATELIMITER_IP_SECONDS_BLOCKED is required")
	}

	return nil
}

var message409 = "you have reached the maximum number of requests or actions allowed within a certain time frame"

func (l *RateLimiter) tokenValidate(tokenString string) error {

	if tokenString == "" {
		return errors.New("authorization token not found")
	}

	_, e := jwtauth.VerifyToken(l.jwtauth, tokenString)

	if e != nil {
		return errors.New("invalid authorization header")
	}

	return nil
}

func (l *RateLimiter) limiter(ctx context.Context, prefixKey string, maxRequests int, secondsBlocked int) *CustomError {

	counterKey := fmt.Sprintf("%s:%d", prefixKey, time.Now().Second())
	blockedKey := fmt.Sprintf("%s:blocked", prefixKey)

	blocked := l.db.Exists(ctx, blockedKey).Val()

	if blocked == 1 {
		return &CustomError{HttpStatus: 429, Message: message409}
	}

	count, _ := l.db.Get(ctx, counterKey).Int()

	count++

	if count > maxRequests {
		l.db.Set(ctx, blockedKey, nil, time.Duration(secondsBlocked)*time.Second)
		return &CustomError{HttpStatus: 429, Message: message409}
	}

	l.db.Set(ctx, counterKey, count, time.Second)

	return nil
}

func (l *RateLimiter) limiterByToken(r *http.Request) *CustomError {

	tokenString := r.Header.Get("API_KEY")

	log.Println(tokenString)

	e := l.tokenValidate(tokenString)

	if e != nil {
		return &CustomError{HttpStatus: 400, Message: e.Error()}
	}

	token, _ := l.jwtauth.Decode(tokenString)

	user, _ := token.Get("user")

	ratelimiterMaxRequests, _ := token.Get("rl-max-requests")
	ratelimiterSecondsBlocked, _ := token.Get("rl-seconds-blocked")

	return l.limiter(r.Context(), user.(string), int(ratelimiterMaxRequests.(float64)), int(ratelimiterSecondsBlocked.(float64)))
}

func (l *RateLimiter) limiterByIP(r *http.Request) *CustomError {

	ip := strings.Split(r.RemoteAddr, ":")[0]

	return l.limiter(r.Context(), ip, l.ratelimiterIpMaxRequests, l.ratelimiterIpSecondsBlocked)
}

func (l *RateLimiter) Limiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var e *CustomError

		if l.ratelimiterTokenOn {
			e = l.limiterByToken(r)
		} else if l.ratelimiterIpOn {
			e = l.limiterByIP(r)
		}

		if e != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(e.HttpStatus)
			w.Write([]byte(`{"message" : "` + e.Error() + `"}`))
			return
		}

		next.ServeHTTP(w, r)
	})
}
