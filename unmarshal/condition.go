package unmarshal

import (
	"fmt"
	"reflect"
	"strings"
)

var (
	_ Condition = conditionRequired(false)
	_ Condition = (*conditionRequireIf)(nil)
	_ Condition = (conditionEnum)(nil)
)

type Condition interface {
	Validate(string) error
}

type MultiCondition interface {
	Condition
	And(Condition) Condition
}

func Required() Condition {
	return conditionRequired(true)
}

func Optional() Condition {
	return conditionRequired(false)
}

type conditionRequired bool

func (c conditionRequired) Validate(s string) error {
	if s != "" || !c {
		return nil
	}
	return ErrEnvNotSet
}

func RequireIf(target interface{}, val string, u Unmarshaler) MultiCondition {
	if val == "" {
		panic("reqif: empty value")
	}

	t := reflect.TypeOf(target)
	if t.Kind() != reflect.Ptr {
		panic(fmt.Errorf("invalit target reference %T", target))
	}

	v := reflect.New(t.Elem())
	if err := u.SetTarget(v.Interface()); err != nil {
		panic(fmt.Errorf("%T: %w", v.Interface(), err))
	}

	if err := u.Unmarshal(val); err != nil {
		panic(err)
	}

	return &conditionRequireIf{
		target:      target,
		value:       v.Elem().Interface(),
		unmarshaler: u,
	}
}

type conditionRequireIf struct {
	target      interface{}
	value       interface{}
	unmarshaler Unmarshaler
	second      Condition
}

func (c *conditionRequireIf) Validate(s string) error {
	if c.required() {
		if c.second != nil {
			return c.second.Validate(s)
		}
		// only check if no second condition
		if s == "" {
			return ErrEnvNotSet
		}
	}
	return nil
}

func (c *conditionRequireIf) And(cond Condition) Condition {
	c.second = cond
	return c
}

func (c *conditionRequireIf) required() bool {
	actual := c.normalize(reflect.ValueOf(c.target).Elem().Interface())
	expected := c.normalize(c.value)

	if actual == expected {
		return true
	}

	switch ev := expected.(type) {
	case int64:
		switch av := actual.(type) {
		case uint64:
			return int64(av) == ev
		case float64:
			return int64(av) == ev
		}
	case uint64:
		switch av := actual.(type) {
		case int64:
			return uint64(av) == ev
		case float64:
			return uint64(av) == ev
		}
	case float64:
		switch av := actual.(type) {
		case int64:
			return float64(av) == ev
		case uint64:
			return float64(av) == ev
		}
	}

	return false
}

func (*conditionRequireIf) normalize(v interface{}) interface{} {
	switch vv := v.(type) {
	case int:
		return int64(vv)
	case int8:
		return int64(vv)
	case int16:
		return int64(vv)
	case int32:
		return int64(vv)
	case uint:
		return uint64(vv)
	case uint8:
		return uint64(vv)
	case uint16:
		return uint64(vv)
	case uint32:
		return uint64(vv)
	case float32:
		return float64(vv)
	default:
		return vv
	}
}

func Enum(s ...string) Condition {
	return conditionEnum(s)
}

type conditionEnum []string

func (c conditionEnum) Validate(s string) error {
	if s == "" {
		return fmt.Errorf("one of [%s] required", strings.Join(c, ","))
	}

	for _, v := range c {
		if s == v {
			return nil
		}
	}
	return fmt.Errorf("string %q not in enum [%s]", s, strings.Join(c, ","))
}
