package reflector

import (
	"fmt"
	"net/url"
	"reflect"
	"time"

	"github.com/albenik/huenv/unmarshal"
)

var registry *unmarshalersRegistry

func init() { //nolint:gochecknoinits
	r := &unmarshalersRegistry{
		unmarshalers: make(map[unmarshal.UnmarshalerName]reflect.Type),
		defaults:     make(map[reflect.Type]unmarshal.UnmarshalerName),
	}
	must(r.RegisterForType("", (*unmarshal.String)(nil)))
	must(r.RegisterForType(false, (*unmarshal.Bool)(nil)))
	must(r.RegisterForType(0, (*unmarshal.Int)(nil)))
	must(r.RegisterForType(int8(0), (*unmarshal.Int8)(nil)))
	must(r.RegisterForType(int16(0), (*unmarshal.Int16)(nil)))
	must(r.RegisterForType(int32(0), (*unmarshal.Int32)(nil)))
	must(r.RegisterForType(int64(0), (*unmarshal.Int64)(nil)))
	must(r.RegisterForType(uint(0), (*unmarshal.Uint)(nil)))
	must(r.RegisterForType(uint8(0), (*unmarshal.Uint8)(nil)))
	must(r.RegisterForType(uint16(0), (*unmarshal.Uint16)(nil)))
	must(r.RegisterForType(uint32(0), (*unmarshal.Uint32)(nil)))
	must(r.RegisterForType(uint64(0), (*unmarshal.Uint64)(nil)))
	must(r.RegisterForType(float32(0), (*unmarshal.Float32)(nil)))
	must(r.RegisterForType(float64(0), (*unmarshal.Float64)(nil)))
	must(r.RegisterForType(time.Time{}, (*unmarshal.Time)(nil)))
	must(r.RegisterForType(time.Duration(0), (*unmarshal.Duration)(nil)))
	must(r.RegisterForType((*url.URL)(nil), (*unmarshal.URL)(nil)))
	must(r.RegisterForType(([]byte)(nil), (*unmarshal.Bytes)(nil)))
	must(r.RegisterForType(([]string)(nil), (*unmarshal.StringsSlice)(nil)))
	registry = r
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

type unmarshalersRegistry struct {
	unmarshalers map[unmarshal.UnmarshalerName]reflect.Type
	defaults     map[reflect.Type]unmarshal.UnmarshalerName
}

func (r *unmarshalersRegistry) Register(u unmarshal.Unmarshaler) error {
	_, err := r.register(u, true)
	return err
}

func (r *unmarshalersRegistry) RegisterForType(i interface{}, u unmarshal.Unmarshaler) error {
	name, err := r.register(u, false)
	if err != nil {
		return err
	}

	r.defaults[reflect.TypeOf(i)] = name
	return nil
}

func (r *unmarshalersRegistry) register(u unmarshal.Unmarshaler, strict bool) (unmarshal.UnmarshalerName, error) {
	utype := reflect.TypeOf(u)
	t := utype
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	name := unmarshal.UnmarshalerName{Package: t.PkgPath(), Type: t.Name()}

	if _, exists := r.unmarshalers[name]; exists {
		if strict {
			return unmarshal.UnmarshalerName{}, newUnmarshalerError(name, ErrAlreadyRegisterd)
		}
		return name, nil
	}

	r.unmarshalers[name] = utype
	return name, nil
}

func (r unmarshalersRegistry) GetUnmarshaler(name unmarshal.UnmarshalerName) (unmarshal.Unmarshaler, error) {
	if ut, ok := r.unmarshalers[name]; ok {
		if um, ok := newValueOfType(ut).Interface().(unmarshal.Unmarshaler); ok {
			return um, nil
		}
		return nil, fmt.Errorf("not an unmarshaler type %s", ut)
	}
	return nil, fmt.Errorf("unmarshaler %q is not registered", name)
}

func (r unmarshalersRegistry) FindUnmarshalerNameForType(t reflect.Type) (unmarshal.UnmarshalerName, bool) {
	if name, ok := r.defaults[t]; ok {
		return name, true
	}
	return unmarshal.UnmarshalerName{}, false
}

func Register(u unmarshal.Unmarshaler) error {
	return registry.Register(u)
}

func RegisterForType(i interface{}, u unmarshal.Unmarshaler) error {
	return registry.RegisterForType(i, u)
}
