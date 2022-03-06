package unmarshal

import (
	"strings"
)

type StringsSlice struct {
	target *[]string
}

func (u *StringsSlice) SetTarget(i interface{}) error {
	if v, ok := i.(*[]string); ok {
		u.target = v
		return nil
	}
	return ErrTypeMismatch
}

func (u *StringsSlice) Unmarshal(s string) error {
	*u.target = strings.Split(s, ",")
	return nil
}
