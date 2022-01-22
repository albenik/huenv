package huenv

import (
	"errors"
	"fmt"
)

var (
	ErrNotGenerated = errors.New("config code mut be generated before use")
	ErrOutdated     = errors.New("generated config unmarshal code is outdated")
)

type KeyError struct {
	key string
	err error
}

func (e *KeyError) Error() string {
	return fmt.Sprintf("env %s: %s", e.key, e.err)
}

func (e *KeyError) Unwrap() error {
	return e.err
}
