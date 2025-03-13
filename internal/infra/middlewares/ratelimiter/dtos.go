package middlewares_ratelimiter

type CustomError struct {
	HttpStatus int
	Message    string
}

func (e *CustomError) Error() string {
	return e.Message
}

type Parameters struct {
	MaxRequests    int
	BlockedSeconds int
}

type Key string
type HttpStatus int
