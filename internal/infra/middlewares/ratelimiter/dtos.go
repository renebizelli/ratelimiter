package middlewares_ratelimiter

type Response struct {
	HttpStatus int
	Message    string
}

func (e *Response) Error() string {
	return e.Message
}

func (e *Response) Ok() bool {
	return e.HttpStatus == 200
}

type Parameters struct {
	MaxRequests    int
	BlockedSeconds int
}

type Key string
