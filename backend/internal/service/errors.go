package service

import "errors"

var ErrNotFound = errors.New("not found")
var ErrInvalidReference = errors.New("referenced resource does not exist")
var ErrNoUpdateFields = errors.New("no fields to update")
var ErrInvalidUserDetails = errors.New("invalid user details")
var ErrInvalidDerivedNameFormat = errors.New("invalid derived name format")
var ErrDerivedNamePropertyInUse = errors.New("property is used by a derived name format")
var ErrInvalidItemTypePropertyOrder = errors.New("invalid item type property order")
var ErrItemReservedByApproval = errors.New("item already approved to another user")
var ErrLocationCycleDetected = errors.New("cycles in location ancestry are not allowed")
var ErrIncorrectPassword = errors.New("incorrect password")
var ErrCannotSetOwnPassword = errors.New("cannot set your own password through this endpoint")
