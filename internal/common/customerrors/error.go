package customerrors

import "fmt"

type HTTPError struct {
	Code    int
	Message string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s", e.Code, e.Message)
}

func Create(code int, message string, args ...any) *HTTPError {
	return &HTTPError{
		Code:    code,
		Message: fmt.Sprintf(message, args...),
	}
}
