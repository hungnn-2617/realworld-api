package middleware

import "errors"

var (
	ErrMissingAuthHeader = errors.New("authorization header is required")
	ErrInvalidAuthFormat = errors.New("authorization header format must be 'Token {token}'")
	ErrEmptyToken        = errors.New("token is empty")
	ErrInvalidToken      = errors.New("invalid or malformed token")
	ErrTokenExpired      = errors.New("token has expired")
)
