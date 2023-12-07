package pmerror

import "errors"

type PMError error

var (
	ErrInvalidInput PMError = errors.New("invalid input")
	ErrNotFound     PMError = errors.New("not found")
	ErrForbidden    PMError = errors.New("forbidden")
	ErrInternal     PMError = errors.New("internal server error")
)
