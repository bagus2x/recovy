package app

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrBadRequest          = errors.New("bad request")
	ErrWrongCredentials    = errors.New("wrong credentials")
	ErrNotFound            = errors.New("not found")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrForbidden           = errors.New("forbidden")
	ErrPermissionDenied    = errors.New("permission denied")
	ErrNotRequiredFields   = errors.New("no such required fields")
	ErrBadQueryParams      = errors.New("invalid query params")
	ErrInternalServerError = errors.New("internal server error")
	ErrRequestTimeoutError = errors.New("request timeout")
)

type RestError struct {
	Code        error  `json:"code"`
	Description string `json:"description"`
	Causes      error  `json:"-"`
}

type RestErrorStatus struct {
	RestError
	Status int `json:"-"`
}

func (e RestError) Error() string {
	return fmt.Sprintf("<code: %s - desc: %s - causes: %s>", e.Code, e.Description, e.Causes)
}

func NewRestError(code error, causes error) error {
	return RestError{
		Code:        code,
		Description: code.Error(),
		Causes:      causes,
	}
}

func NewRestErrorWithMessage(code error, desc string, causes error) error {
	return RestError{
		Code:        code,
		Description: desc,
		Causes:      causes,
	}
}

func ParseErrors(err error) RestErrorStatus {
	if e, ok := err.(RestError); ok {
		return RestErrorStatus{
			RestError: e,
			Status:    statusCode(e.Code),
		}
	}

	return RestErrorStatus{
		RestError: RestError{
			Code:        ErrInternalServerError,
			Description: ErrInternalServerError.Error(),
			Causes:      err,
		},
		Status: 500,
	}
}

func statusCode(err error) int {
	switch err {
	case ErrBadRequest:
		return http.StatusBadRequest
	}

	return http.StatusInternalServerError
}
