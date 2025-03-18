package middlewares_ratelimiter

import (
	"net/http"
	"strings"
)

type BasedOnIP struct {
	core           CoreInterface
	headerByStuffs *HeaderByStuffs
	on             bool
	parameters     *Parameters
}

func NewBasedOnIP(core CoreInterface, headerByStuffs *HeaderByStuffs, on bool, maxRequests, blockedSeconds int) *BasedOnIP {
	return &BasedOnIP{
		core:           core,
		headerByStuffs: headerByStuffs,
		on:             on,
		parameters: &Parameters{
			MaxRequests:    maxRequests,
			BlockedSeconds: blockedSeconds,
		},
	}
}

func (l *BasedOnIP) validate() *Response {

	if l.parameters.MaxRequests == 0 {
		return &Response{Message: "RATELIMITER_IP_MAX_REQUESTS is required", HttpStatus: http.StatusBadRequest}
	} else if l.parameters.BlockedSeconds == 0 {
		return &Response{Message: "RATELIMITER_IP_BLOCKED_SECONDS is required", HttpStatus: http.StatusBadRequest}
	}

	return nil
}

func (l *BasedOnIP) Limiter(r *http.Request) *Response {

	if !l.on {
		return nil
	}

	if l.headerByStuffs.IsByPass(r) {
		return nil
	}

	if e := l.validate(); e != nil {
		return e
	}

	ip := strings.Split(r.RemoteAddr, ":")[0]
	key := Key(ip)

	parameters := &Parameters{
		MaxRequests:    l.parameters.MaxRequests,
		BlockedSeconds: l.parameters.BlockedSeconds,
	}

	if status := l.core.Limiter(r.Context(), key, parameters); status != http.StatusOK {
		return &Response{HttpStatus: status}
	}

	return nil
}
