package customerrors

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)
