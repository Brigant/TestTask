package model

import "errors"

var (
	ErrNotFound             = errors.New("nothing was found")
	ErrWrongBody            = errors.New("wrong body")
	ErrInvalidSigningMethod = errors.New("invalid signing metod")
	ErrWrongTokenClaimType  = errors.New("token claims are not of type *tokenClaims")
	ErrTimeout              = errors.New("timeout")
	ErrDatabaseViolation    = errors.New("database violtion")
	ErrNothinChanged        = errors.New("nothing has changed")
)
