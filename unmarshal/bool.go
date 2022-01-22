package unmarshal

import (
	"strconv"
)

type Bool struct {
	target *bool
}

func (u *Bool) SetTarget(i interface{}) error {
	if v, ok := i.(*bool); ok {
		u.target = v
		return nil
	}
	return ErrTypeMismatch
}

func (u *Bool) Unmarshal(s string) (err error) {
	*u.target, err = strconv.ParseBool(s)
	return
}
