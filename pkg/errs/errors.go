package errs

import (
	"errors"
)

var (
	FioFailedErr    = errors.New("неправильные данные")
	UserNotFoundErr = errors.New("пользователь не найден")
)
