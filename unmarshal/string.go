package unmarshal

type String struct {
	target *string
}

func (u *String) SetTarget(i interface{}) error {
	if v, ok := i.(*string); ok {
		u.target = v
		return nil
	}
	return ErrTypeMismatch
}

func (u *String) Unmarshal(s string) error {
	*u.target = s
	return nil
}
