package domain

import "errors"

// Domain-level error sentinels.
var (
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrSIAKADUnavailable   = errors.New("SIAKAD service unavailable")
	ErrForbidden           = errors.New("forbidden")
	ErrNotFound            = errors.New("not found")
	ErrDuplicateAlias      = errors.New("alias already taken")
	ErrInvalidSignature    = errors.New("invalid webhook signature")
	ErrTokenInvalid        = errors.New("token is invalid")
	ErrTokenExpired        = errors.New("token has expired")
	ErrDeviceRejected      = errors.New("device rejected by SIAKAD")
	ErrSessionInitFailed   = errors.New("session initialization failed")
	ErrInsufficientScope   = errors.New("insufficient scope")
)
