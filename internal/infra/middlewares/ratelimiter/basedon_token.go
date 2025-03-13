package middlewares_ratelimiter

import (
	"errors"
	"net/http"

	"github.com/go-chi/jwtauth"
	"github.com/lestrrat-go/jwx/jwt"
)

type BasedOnToken struct {
	jwtauth  *jwtauth.JWTAuth
	_default *Parameters
}

func NewBasedOnToken(jwtauth *jwtauth.JWTAuth, defaultMaxRequests, defaultBlockedSeconds int) *BasedOnToken {
	return &BasedOnToken{
		jwtauth:  jwtauth,
		_default: &Parameters{MaxRequests: defaultMaxRequests, BlockedSeconds: defaultBlockedSeconds},
	}
}

func (l *BasedOnToken) Validate(r *http.Request) error {

	tokenString := r.Header.Get("API_KEY")

	if tokenString == "" {
		return errors.New("authorization token not found")
	}

	_, e := jwtauth.VerifyToken(l.jwtauth, tokenString)

	if e != nil {
		return errors.New("invalid authorization header")
	}

	return nil
}

func (l *BasedOnToken) Parse(r *http.Request) (Key, *Parameters) {

	tokenString := r.Header.Get("API_KEY")

	token, _ := l.jwtauth.Decode(tokenString)

	key := l.extractKey(token)

	maxRequests := l.extractIntValue(token, "rl-max-requests", l._default.MaxRequests)
	blockedSeconds := l.extractIntValue(token, "rl-seconds-blocked", l._default.BlockedSeconds)

	return key, &Parameters{
		MaxRequests:    maxRequests,
		BlockedSeconds: blockedSeconds,
	}
}

func (l *BasedOnToken) extractKey(token jwt.Token) Key {
	user, _ := token.Get("user")
	key := Key(user.(string))
	return key
}

func (l *BasedOnToken) extractIntValue(token jwt.Token, claimName string, defaultValue int) int {
	value := defaultValue
	if v, e := token.Get(claimName); e {
		value = int(v.(float64))
	}
	return value
}
