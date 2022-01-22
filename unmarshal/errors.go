package unmarshal

import (
	"errors"
)

var (
	ErrEnvNotSet      = errors.New("environment variable not set or empty")
	ErrInvalidOptions = errors.New("invalid options")
	ErrTypeMismatch   = errors.New("value type mismatch")
	ErrOverflow       = errors.New("value overflow")
)
