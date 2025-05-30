package pkg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// Application error codes.
// These codes are used to represent and pass specific application errors in the project.
const (
	ECONFLICT            = "conflict"
	EINTERNAL            = "internal"
	EINVALID             = "invalid"
	ENOTFOUND            = "not_found"
	ENOTIMPLEMENTED      = "not_implemented"
	EUNAUTHORIZED        = "unauthorized"
	EUNPROCESSABLEENTITY = "unprocessable_entity"
	ETOOMANYREQUESTS     = "too_many_requests"
	EFORBIDDEN           = "forbidden"
	ECANCELED            = "canceled"
)

// Lookup of application error codes to HTTP status codes.
var codes = map[string]int{
	ECONFLICT:            http.StatusConflict,
	EINVALID:             http.StatusBadRequest,
	ENOTFOUND:            http.StatusNotFound,
	ENOTIMPLEMENTED:      http.StatusNotImplemented,
	EUNAUTHORIZED:        http.StatusUnauthorized,
	EUNPROCESSABLEENTITY: http.StatusUnprocessableEntity,
	EINTERNAL:            http.StatusInternalServerError,
	ETOOMANYREQUESTS:     http.StatusTooManyRequests,
	EFORBIDDEN:           http.StatusForbidden,
	ECANCELED:            499,
}

// Error represents a structured application error.
type Error struct {
	// Optional wrapped error.
	Err error `json:"error,omitempty"`

	// Machine-readable error code.
	Code string `json:"code"`

	// Human-readable error message.
	Message string `json:"message"`
}

// Error() implements the error interface.
func (e *Error) Error() string {
	jsonErr, err := json.Marshal(e)

	if err != nil {
		errStrg := fmt.Sprintf("code:%s message:%s", e.Code, e.Message)
		if e.Err != nil {
			errStrg += fmt.Sprintf(" error:%s", e.Err)
		}
		return errStrg
	}

	return string(jsonErr)

}

// Is implements the error comparison based on the error code.
func (e *Error) Is(err error) bool {
	var eErr *Error
	if errors.As(err, &eErr) {
		return eErr.Code == e.Code
	}
	return false
}

// ErrorCode unwraps an application error and returns its code.
// Non-application errors always return EINTERNAL.
func ErrorCode(err error) string {
	var e *Error
	if err == nil {
		return ""
	} else if errors.As(err, &e) {
		return e.Code
	} else if errors.Is(err, context.Canceled) {
		return ECANCELED
	}
	return EINTERNAL
}

// ErrorMessage unwraps an application error and returns its message.
// Non-application errors always return "Internal error".
func ErrorMessage(err error) string {
	if err == nil {
		return ""
	}

	var e *Error
	if errors.As(err, &e) {
		return e.Message
	}

	return "Internal error."
}

// Errorf is a helper function to return an WriteError with a given code and formatted message.
func Errorf(code, msg string, rawErr error) *Error {
	return &Error{
		Code:    code,
		Message: msg,
		Err:     rawErr,
	}
}

// ErrorStatusCode returns the associated HTTP status code for an arh.WriteError code.
func ErrorStatusCode(code string) int {
	if v, ok := codes[code]; ok {
		return v
	}
	return http.StatusInternalServerError
}

// FromErrorStatusCode returns the associated arh.code for a HTTP status code.
func FromErrorStatusCode(code int) string {
	for k, v := range codes {
		if v == code {
			return k
		}
	}
	return EINTERNAL
}
