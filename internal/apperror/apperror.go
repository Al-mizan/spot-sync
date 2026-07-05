package apperror

import "net/http"

// AppError is a domain error that carries an HTTP status code.
type AppError struct {
	Code    int    // HTTP status code
	Message string // user-facing message
	Err     error  // wrapped original error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error { return e.Err }

func NewNotFound(err error, msg string) *AppError {
	return &AppError{Code: http.StatusNotFound, Message: msg, Err: err}
}

func NewConflict(err error, msg string) *AppError {
	return &AppError{Code: http.StatusConflict, Message: msg, Err: err}
}

func NewUnauthorized(err error, msg string) *AppError {
	return &AppError{Code: http.StatusUnauthorized, Message: msg, Err: err}
}

func NewForbidden(err error, msg string) *AppError {
	return &AppError{Code: http.StatusForbidden, Message: msg, Err: err}
}

func NewInternal(err error, msg string) *AppError {
	return &AppError{Code: http.StatusInternalServerError, Message: msg, Err: err}
}

func NewBadRequest(err error, msg string) *AppError {
	return &AppError{Code: http.StatusBadRequest, Message: msg, Err: err}
}
