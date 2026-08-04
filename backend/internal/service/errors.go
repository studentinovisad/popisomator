package service

import "errors"

var ErrNotFound = errors.New("not found")
var ErrInvalidReference = errors.New("referenced resource does not exist")
