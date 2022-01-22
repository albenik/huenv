package unmarshal

import (
	"strconv"
)

const intSize = 32 << (^uint(0) >> 63) // 32 or 64

type Int struct {
	target *int
}

func (u *Int) SetTarget(i interface{}) error {
	if v, ok := i.(*int); ok {
		u.target = v
		return nil
	}
	return ErrTypeMismatch
}

func (u *Int) Unmarshal(s string) error {
	val, err := strconv.ParseInt(s, 10, intSize)
	if err != nil {
		return err
	}
	*u.target = int(val)
	return nil
}

type Int8 struct {
	target *int8
}

func (u *Int8) SetTarget(i interface{}) error {
	if v, ok := i.(*int8); ok {
		u.target = v
		return nil
	}
	return ErrTypeMismatch
}

func (u *Int8) Unmarshal(str string) error {
	val, err := strconv.ParseInt(str, 10, 8)
	if err != nil {
		return err
	}
	*u.target = int8(val)
	return nil
}

type Int16 struct {
	target *int16
}

func (u *Int16) SetTarget(i interface{}) error {
	if v, ok := i.(*int16); ok {
		u.target = v
		return nil
	}
	return ErrTypeMismatch
}

func (u *Int16) Unmarshal(str string) error {
	val, err := strconv.ParseInt(str, 10, 16)
	if err != nil {
		return err
	}
	*u.target = int16(val)
	return nil
}

type Int32 struct {
	target *int32
}

func (u *Int32) SetTarget(i interface{}) error {
	if v, ok := i.(*int32); ok {
		u.target = v
		return nil
	}
	return ErrTypeMismatch
}

func (u *Int32) Unmarshal(str string) error {
	val, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		return err
	}
	*u.target = int32(val)
	return nil
}

type Int64 struct {
	target *int64
}

func (u *Int64) SetTarget(i interface{}) error {
	if v, ok := i.(*int64); ok {
		u.target = v
		return nil
	}
	return ErrTypeMismatch
}

func (u *Int64) Unmarshal(str string) error {
	val, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return err
	}
	*u.target = val
	return nil
}

type Uint struct {
	target *uint
}

func (u *Uint) SetTarget(i interface{}) error {
	if v, ok := i.(*uint); ok {
		u.target = v
		return nil
	}
	return ErrTypeMismatch
}

func (u *Uint) Unmarshal(str string) error {
	val, err := strconv.ParseUint(str, 10, intSize)
	if err != nil {
		return err
	}
	*u.target = uint(val)
	return nil
}

type Uint8 struct {
	target *uint8
}

func (u *Uint8) SetTarget(i interface{}) error {
	if v, ok := i.(*uint8); ok {
		u.target = v
		return nil
	}
	return ErrTypeMismatch
}

func (u *Uint8) Unmarshal(str string) error {
	val, err := strconv.ParseUint(str, 10, 8)
	if err != nil {
		return err
	}
	*u.target = uint8(val)
	return nil
}

type Uint16 struct {
	target *uint16
}

func (u *Uint16) SetTarget(i interface{}) error {
	if v, ok := i.(*uint16); ok {
		u.target = v
		return nil
	}
	return ErrTypeMismatch
}

func (u *Uint16) Unmarshal(str string) error {
	val, err := strconv.ParseUint(str, 10, 16)
	if err != nil {
		return err
	}
	*u.target = uint16(val)
	return nil
}

type Uint32 struct {
	target *uint32
}

func (u *Uint32) SetTarget(i interface{}) error {
	if v, ok := i.(*uint32); ok {
		u.target = v
		return nil
	}
	return ErrTypeMismatch
}

func (u *Uint32) Unmarshal(str string) error {
	val, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		return err
	}
	*u.target = uint32(val)
	return nil
}

type Uint64 struct {
	target *uint64
}

func (u *Uint64) SetTarget(i interface{}) error {
	if v, ok := i.(*uint64); ok {
		u.target = v
		return nil
	}
	return ErrTypeMismatch
}

func (u *Uint64) Unmarshal(str string) error {
	val, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return err
	}
	*u.target = val
	return nil
}
