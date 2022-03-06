package unmarshal

import (
	"encoding/base64"
)

type Bytes struct {
	target *[]byte
}

func (u *Bytes) SetTarget(t interface{}) error {
	if b, ok := t.(*[]byte); ok {
		u.target = b
		return nil
	}
	return ErrTypeMismatch
}

func (u *Bytes) Unmarshal(str string) (err error) {
	*u.target, err = base64.StdEncoding.DecodeString(str)
	return
}
