package unmarshal

import (
	"fmt"
)

type Unmarshaler interface {
	SetTarget(interface{}) error
	Unmarshal(string) error
}

type UnmarshalerName struct {
	Package string // package import path or name
	Type    string // type name without package
}

func (n UnmarshalerName) String() string {
	return fmt.Sprintf("%s %s", n.Package, n.Type)
}

type Target struct {
	unmarshaler Unmarshaler
	condition   Condition
}

func (t *Target) Unmarshal(s string) error {
	if err := t.condition.Validate(s); err != nil {
		return err
	}
	if s == "" {
		return nil
	}
	return t.unmarshaler.Unmarshal(s)
}

func (t *Target) Dependency() interface{} {
	if c, ok := t.condition.(*conditionRequireIf); ok {
		return c.target
	}
	return nil
}

func NewTarget(t interface{}, u Unmarshaler, c Condition) *Target {
	if err := u.SetTarget(t); err != nil {
		panic(err)
	}

	return &Target{
		unmarshaler: u,
		condition:   c,
	}
}
