package svc

import (
	"fmt"
)

type Error interface {
	Code() int
	Message() string
}

type err struct {
	code    int
	err     error
	message string
}

const (
	ErrCodeSystem = iota + 1
	ErrCodeUnauthorized
	ErrCodePermissionDenied
	ErrCodeInvalidParams
)

func (e *err) Error() string {
	return fmt.Sprint(e.err)
}

func (e *err) Code() int {
	return e.code
}

func (e *err) Message() string {
	return e.message
}

func NewError(e error, code int, msg string) error {
	return &err{
		code:    code,
		err:     e,
		message: msg,
	}
}

func ErrSystem(e error, msg ...string) error {
	m := "Something went wrong, please try again later."
	if len(msg) > 0 {
		m = msg[0]
	}
	return NewError(e, ErrCodeSystem, m)
}

func ErrInvalidParams(e error, msg ...string) error {
	m := "Invalid params."
	if len(msg) > 0 {
		m = msg[0]
	}
	return NewError(e, ErrCodeInvalidParams, m)
}

func ErrPermissionDenied(e error, msg ...string) error {
	m := "Permission denied."
	if len(msg) > 0 {
		m = msg[0]
	}
	return NewError(e, ErrCodePermissionDenied, m)
}

func ErrUnauthorized(e error, msg ...string) error {
	m := "Unauthorized."
	if len(msg) > 0 {
		m = msg[0]
	}
	return NewError(e, ErrCodeUnauthorized, m)
}
