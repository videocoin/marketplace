package accounts

import "errors"

var (
	ErrInvalidAddress           = errors.New("invalid address")
	ErrAddressAlreadyRegistered = errors.New("address already registered")
	ErrAddressNotRegistered     = errors.New("address not registered")
	ErrInvalidImageData         = errors.New("invalid image data")
)
