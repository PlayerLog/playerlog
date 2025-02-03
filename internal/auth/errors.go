package auth

import "errors"

var (
	ErrUserAlreadyExist     = errors.New("user already exists.")
	ErrNoUserFound          = errors.New("no user found.")
	ErrWrongEmailOrPassword = errors.New("wrong email or password.")
)
