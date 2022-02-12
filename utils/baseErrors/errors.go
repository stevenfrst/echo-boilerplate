package baseErrors

import "errors"

var (
	ErrUserEmailUsed       = errors.New("email already used")
	ErrUserInactive        = errors.New("user not verified")
	ErrPasswordNotMatch    = errors.New("password does not match")
	ErrOldPasswordNotMatch = errors.New("old password does not match")
	ErrDeleteAccount       = errors.New("errors when deleting account")
	ErrWrongQueryParams    = errors.New("wrong params")
	ErrInvalidPayload      = errors.New("invalid payload")
	ErrUserEmailNotFound   = errors.New("email not found")
)
