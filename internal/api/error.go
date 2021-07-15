package api

import "errors"

var (
	ErrInvalidAddress           = errors.New("invalid address")
	ErrAddressAlreadyRegistered = errors.New("address already registered")
	ErrAddressNotRegistered     = errors.New("address not registered")
	ErrInvalidImageData         = errors.New("invalid image data")
	ErrUsernameAlreadyUsed      = errors.New("username already used")
	ErrInvalidSignature         = errors.New("invalid signature")

	ErrUnsupportedContentType = errors.New("unsupported content type")
	ErrInvalidVideo           = errors.New("invalid video")
	ErrInvalidMedia           = errors.New("invalid media")

	ErrInvalidEncPublicKey = errors.New("invalid encryption public key")
)
