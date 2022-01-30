package unmarshal

import (
	"errors"
)

var (
	ErrEnvNotSet    = errors.New("variable not set or empty")
	ErrTypeMismatch = errors.New("value type mismatch")
)
