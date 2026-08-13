package service

import "errors"

var ErrNotFound = errors.New("not found")
var ErrInvalidReference = errors.New("referenced resource does not exist")
var ErrInvalidRegistrationStatus = errors.New("user registration is no longer pending")
