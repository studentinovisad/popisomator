package service

import "errors"

var ErrNotFound = errors.New("not found")
var ErrInvalidReference = errors.New("referenced resource does not exist")
var ErrNoUpdateFields = errors.New("no fields to update")
