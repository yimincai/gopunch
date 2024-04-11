package errs

import "errors"

var (
	ErrInternalError = errors.New("internal error")
	ErrLoginFailed   = errors.New("login failed")
	ErrPunchFailed   = errors.New("punch failed")
	ErrUserNotFound  = errors.New("user not found")
	ErrUserDisabled  = errors.New("user disabled")
)
