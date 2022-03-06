package unmarshal

import (
	"strconv"
)

type Float32 struct {
	target *float32
}

func (u *Float32) SetTarget(i interface{}) error {
	if v, ok := i.(*float32); ok {
		u.target = v
		return nil
	}
	return ErrTypeMismatch
}

func (u *Float32) Unmarshal(str string) error {
	v, err := strconv.ParseFloat(str, 32)
	if err != nil {
		return err
	}

	*u.target = float32(v)
	return nil
}

type Float64 struct {
	target *float64
}

func (u *Float64) SetTarget(i interface{}) error {
	if v, ok := i.(*float64); ok {
		u.target = v
		return nil
	}
	return ErrTypeMismatch
}

func (u *Float64) Unmarshal(str string) (err error) {
	*u.target, err = strconv.ParseFloat(str, 64)
	return
}
