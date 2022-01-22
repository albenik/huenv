package reflector

import (
	"github.com/albenik/huenv/unmarshal"
)

var _ unmarshal.Unmarshaler = (*TestUnmarshaler1)(nil)

type TestUnmarshaler1 struct{}

func (t TestUnmarshaler1) SetTarget(interface{}) error {
	return nil
}

func (t TestUnmarshaler1) Unmarshal(string) error {
	return nil
}
