package middlewares_ratelimiter

import (
	"net/http"

	pkg_utils "github.com/renebizelli/ratelimiter/pkg/utils"
)

type BasedOnToken struct {
	core           CoreInterface
	jwt            *pkg_utils.Jwt
	headerByStuffs *HeaderByStuffs
	on             bool
	_default       *Parameters
}

func NewBasedOnToken(core CoreInterface, jwt *pkg_utils.Jwt, headerByStuffs *HeaderByStuffs, on bool, defaultMaxRequests, defaultBlockedSeconds int) *BasedOnToken {
	return &BasedOnToken{
		core:           core,
		jwt:            jwt,
		headerByStuffs: headerByStuffs,
		on:             on,
		_default:       &Parameters{MaxRequests: defaultMaxRequests, BlockedSeconds: defaultBlockedSeconds},
	}
}

func (l *BasedOnToken) validate(tokenString string) Response {

	if e := l.jwt.VerifyToken(tokenString); e != nil {
		return Response{Message: "invalid authorization header", HttpStatus: http.StatusBadRequest}
	}

	return Response{HttpStatus: 200}
}

func (l *BasedOnToken) Limiter(r *http.Request, ch chan<- Response) {

	if !l.on {
		ch <- Response{HttpStatus: 200}
		return
	}

	var tokenString string

	if t, e := l.headerByStuffs.GetAPIKey(r); e != nil {
		ch <- Response{HttpStatus: 200}
		return
	} else {
		tokenString = t
	}

	if e := l.validate(tokenString); !e.Ok() {
		ch <- e
		return
	}

	token := l.jwt.Decode(tokenString)

	key := Key(l.jwt.ExtractStringValue(token, "key"))

	maxRequests := l.jwt.ExtractIntValue(token, "rl-max-requests", l._default.MaxRequests)
	blockedSeconds := l.jwt.ExtractIntValue(token, "rl-seconds-blocked", l._default.BlockedSeconds)

	parameters := &Parameters{
		MaxRequests:    maxRequests,
		BlockedSeconds: blockedSeconds,
	}

	if status := l.core.Limiter(r.Context(), key, parameters); status != http.StatusOK {
		ch <- Response{HttpStatus: status}
		return
	}

	l.headerByStuffs.Set(r)

	ch <- Response{HttpStatus: 200}
}
