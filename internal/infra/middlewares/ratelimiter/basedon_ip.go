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

func (l *BasedOnIP) validate() Response {

	if l.parameters.MaxRequests == 0 {
		return Response{Message: "RATELIMITER_IP_MAX_REQUESTS is required", HttpStatus: http.StatusBadRequest}
	} else if l.parameters.BlockedSeconds == 0 {
		return Response{Message: "RATELIMITER_IP_BLOCKED_SECONDS is required", HttpStatus: http.StatusBadRequest}
	}

	return Response{HttpStatus: 200}
}

func (l *BasedOnIP) Limiter(r *http.Request, ch chan<- Response) {

	if !l.on {
		ch <- Response{HttpStatus: 200}
		return
	}

	if l.headerByStuffs.IsByPass(r) {
		ch <- Response{HttpStatus: 200}
		return
	}

	if e := l.validate(); !e.Ok() {
		ch <- e
		return
	}

	ip := strings.Split(r.RemoteAddr, ":")[0]
	key := Key(ip)

	parameters := &Parameters{
		MaxRequests:    l.parameters.MaxRequests,
		BlockedSeconds: l.parameters.BlockedSeconds,
	}

	if status := l.core.Limiter(key, parameters); status != http.StatusOK {
		ch <- Response{HttpStatus: status}
		return
	}

	ch <- Response{HttpStatus: 200}
}
