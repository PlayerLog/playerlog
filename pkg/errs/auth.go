package errs

import "errors"

var (
	UserAlreadyExist = errors.New("User already exists.")
)
