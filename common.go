package pkt

import (
	"errors"
)

var (
	// ErrTooShort packet is too short
	ErrTooShort = errors.New("packet is too short")
	// ErrBadFormat packet format is invalid
	ErrBadFormat = errors.New("packet bad format")
)
