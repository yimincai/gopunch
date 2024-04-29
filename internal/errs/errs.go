package errs

import (
	"errors"

	"github.com/yimincai/gopunch/pkg/logger"
)

var (
	ErrInternalError     = errors.New("internal error, please contact admin 👨‍💻")
	ErrForbidden         = errors.New("forbidden, please contact admin 👨‍💻")
	ErrSendingMessage    = errors.New("error while sending message, 👨‍💻")
	ErrDayOffAlreadySet  = errors.New("day off already set 📅")
	ErrPunchFailed       = errors.New("punch failed 🥊")
	ErrUserNotFound      = errors.New("user not found, please register first 👨‍💻")
	ErrUserNotEnabled    = errors.New("user not enabled 🤕, please contact admin 👨‍💻")
	ErrInvalidDate       = errors.New("invalid date, please check the date format 📅")
	ErrDeleteUserFailed  = errors.New("delete user failed 😣, please contact admin 👨‍💻")
	ErrLoginVerifyFailed = errors.New("login verify failed 😦, please check your account and password 🤔")
	ErrLoginFailed       = errors.New("login failed 😦, please check your account and password and force register again 🤔")
	ErrPunchOnTimeFailed = errors.New("punch on time failed 🕰️")
)

func LogError(err error) {
	if err != nil {
		logger.Error(err)
	}
}
