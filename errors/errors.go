package errors

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	StatusCode int    `json:"statusCode"`
	StatusType string `json:"statusType"`
	Message    string `json:"message"`
}

func newError(statusCode int, err error) *Error {
	return &Error{statusCode, http.StatusText(statusCode), err.Error()}
}

func (e *Error) Error() string {
	js, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		panic("unreachable")
	}
	return string(js)
}

func WrapError(err error) *Error {
	if e, ok := err.(*Error); ok {
		return e
	}
	return newError(http.StatusInternalServerError, err)
}

func UnsupportedMediaType(err error) *Error {
	return newError(http.StatusUnsupportedMediaType, err)
}

func BadRequest(err error) *Error {
	return newError(http.StatusBadRequest, err)
}

func NotFound(err error) *Error {
	return newError(http.StatusNotFound, err)
}

func Conflict(err error) *Error {
	return newError(http.StatusConflict, err)
}
