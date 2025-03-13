package middlewares_ratelimiter

import (
	"errors"
	"net/http"
	"strings"
)

type BasedOnIP struct {
	parameters *Parameters
}

func NewBasedOnIP(maxRequests, blockedSeconds int) *BasedOnIP {
	return &BasedOnIP{
		parameters: &Parameters{
			MaxRequests:    maxRequests,
			BlockedSeconds: blockedSeconds,
		},
	}
}

func (l *BasedOnIP) Validate() error {

	if l.parameters.MaxRequests == 0 {
		return errors.New("RATELIMITER_IP_MAX_REQUESTS is required")
	} else if l.parameters.BlockedSeconds == 0 {
		return errors.New("RATELIMITER_IP_BLOCKED_SECONDS is required")
	}

	return nil
}

func (l *BasedOnIP) Parse(r *http.Request) (Key, *Parameters) {

	ip := strings.Split(r.RemoteAddr, ":")[0]
	key := Key(ip)

	return key, &Parameters{
		MaxRequests:    l.parameters.MaxRequests,
		BlockedSeconds: l.parameters.BlockedSeconds,
	}

}
