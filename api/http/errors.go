package http

type Error struct {
	message string
	code    int
}

func NewApiError(message string, code int) Error {
	return Error{
		message: message,
		code:    code,
	}
}
