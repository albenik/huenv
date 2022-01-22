package reflector

import (
	"errors"
	"fmt"

	"github.com/albenik/huenv/unmarshal"
)

var (
	ErrAlreadyRegisterd    = errors.New("already registered")
	ErrInvalidEnvTag       = errors.New("env tag empty or invalid")
	ErrUnmarshalerNotFound = errors.New("unmarshaler not found")
)

type InvalidConfigTypeError string

func (e InvalidConfigTypeError) Error() string {
	return fmt.Sprintf("provided type %s is not a pointer to the struct", string(e))
}

type UnsupportedFieldTypeError struct {
	Field string
	Type  string
}

func (e *UnsupportedFieldTypeError) Error() string {
	return fmt.Sprintf("unsupported field type %s %s", e.Field, e.Type)
}

func newUnmarshalerError(name unmarshal.UnmarshalerName, err error) error {
	return fmt.Errorf("unmarshaler `%s`: %w", name, err)
}

func newFieldError(name string, err error) error {
	return fmt.Errorf("field %q: %w", name, err)
}
