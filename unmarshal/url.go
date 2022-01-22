package unmarshal

import (
	"net/url"
)

type URL struct {
	target    **url.URL
}

func (u *URL) SetTarget(i interface{}) error {
	if v, ok := i.(**url.URL); ok {
		u.target = v
		return nil
	}
	return ErrTypeMismatch
}

func (u *URL) Unmarshal(str string) (err error) {
	*u.target, err = url.Parse(str)
	return
}
