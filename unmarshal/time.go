package unmarshal

import (
	"time"
)

type Time struct {
	target *time.Time
}

func (u *Time) SetTarget(i interface{}) error {
	if v, ok := i.(*time.Time); ok {
		u.target = v
		return nil
	}
	return ErrTypeMismatch
}

func (u *Time) Unmarshal(str string) (err error) {
	*u.target, err = time.Parse(time.RFC3339, str)
	return
}

type Duration struct {
	target *time.Duration
}

func (u *Duration) SetTarget(i interface{}) error {
	if v, ok := i.(*time.Duration); ok {
		u.target = v
		return nil
	}
	return ErrTypeMismatch
}

func (u *Duration) Unmarshal(str string) (err error) {
	*u.target, err = time.ParseDuration(str)
	return
}
