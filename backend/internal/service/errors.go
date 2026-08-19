package service

import "errors"

var ErrNotFound = errors.New("not found")
var ErrInvalidReference = errors.New("referenced resource does not exist")
var ErrNoUpdateFields = errors.New("no fields to update")
var ErrInvalidDerivedNameFormat = errors.New("invalid derived name format")
var ErrDerivedNamePropertyInUse = errors.New("property is used by a derived name format")
var ErrItemReservedByApproval = errors.New("item already approved to another user")
