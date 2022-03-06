package reflector

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/albenik/huenv/unmarshal"
)

const (
	envTagName         = "env"
	unmarshalerTagName = "unmarshaler"
)

type Reflector struct {
	pkgs   map[string]struct{}
	envs   map[string]*Target
	fields map[string]*TargetField
}

func New() *Reflector {
	return &Reflector{
		pkgs:   make(map[string]struct{}),
		envs:   make(map[string]*Target),
		fields: make(map[string]*TargetField),
	}
}

func (r *Reflector) Reflect(conf interface{}) (*Result, error) {
	t := reflect.TypeOf(conf)
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		return nil, InvalidConfigTypeError(t.String())
	}

	ts := t.Elem()
	vs := reflect.ValueOf(conf).Elem()

	if err := r.processStruct(ts, &vs, "", "", nil); err != nil {
		return nil, err
	}

	tpref, tname := parseTypeName(ts)
	return &Result{
		ConfigPkg:  tpref,
		ConfigType: tname,
		Packages:   r.pkgs,
		Envs:       r.envs,
	}, nil
}

func (r *Reflector) processStruct(
	typ reflect.Type,
	val *reflect.Value,
	envPrefix, fieldPrefix string,
	parentCond *ConditionRequireIf,
) error {
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fvalue := val.Field(i)

		longFieldName := field.Name
		if fieldPrefix != "" {
			longFieldName = fmt.Sprintf("%s.%s", fieldPrefix, longFieldName)
		}

		uname, hasUnmarshaler := findUnmarshalerNameForField(&field)

		if !hasUnmarshaler {
			if processed, err := r.tryProcessSublevel(&field, &fvalue, longFieldName, envPrefix, parentCond); err != nil {
				return newFieldError(longFieldName, err)
			} else if processed {
				continue
			}
		}

		tag, err := parseEnvTag(field.Tag.Get(envTagName))
		if err != nil {
			return newFieldError(longFieldName, err)
		}

		cond, err := r.constructCondition(longFieldName, tag.Cond)
		if err != nil {
			return newFieldError(longFieldName, err)
		}

		if !hasUnmarshaler {
			return newFieldError(longFieldName, ErrUnmarshalerNotFound)
		}

		if _, err = registry.GetUnmarshaler(uname); err != nil {
			return newFieldError(longFieldName, fmt.Errorf("can not create unmarshaler: %w", err))
		}

		envName := envPrefix + tag.EnvName
		if _, ok := r.envs[envName]; ok {
			return newFieldError(longFieldName, fmt.Errorf("environment variable name %q already used", envName))
		}

		info := &Target{
			Field: &TargetField{
				Name:        longFieldName,
				Unmarshaler: uname,
			},
			Condition: combineConditions(parentCond, cond),
		}

		r.pkgs[uname.Package] = struct{}{}
		r.envs[envName] = info
		r.fields[longFieldName] = info.Field
	}

	return nil
}

func (r *Reflector) tryProcessSublevel(
	field *reflect.StructField,
	fval *reflect.Value,
	longFieldName, envPrefix string,
	parentCond *ConditionRequireIf,
) (bool, error) {
	t := dereferenceT(field.Type)
	if t.Kind() != reflect.Struct || field.Tag.Get(unmarshalerTagName) != "" {
		return false, nil
	}

	var condition *ConditionRequireIf
	if s := field.Tag.Get(envTagName); s != "" {
		tag, err := parseEnvTag(s)
		if err != nil {
			return false, err
		}

		if !strings.HasSuffix(tag.EnvName, "*") {
			return false, ErrInvalidEnvTag
		}
		envPrefix += strings.TrimSuffix(tag.EnvName, "*")

		c, err := r.constructCondition(longFieldName, tag.Cond)
		if err != nil {
			return false, err
		}
		condition, _ = c.(*ConditionRequireIf) // this is expected behaviour: specific type or nil
	}

	if condition == nil {
		condition = parentCond
	}

	if field.Type.Kind() == reflect.Ptr {
		fval.Set(newValueOfType(field.Type))
	}

	if err := r.processStruct(t, dereferenceV(fval), envPrefix, longFieldName, condition); err != nil {
		return false, err
	}

	return true, nil
}

func (r *Reflector) constructCondition(fname string, cond interface{}) (interface{}, error) {
	if cond == nil {
		return ConditionRequired(true), nil
	}

	switch c := cond.(type) {
	case tagCondOptional:
		return ConditionRequired(false), nil

	case *tagCondIf:
		tname := c.Field
		if i := strings.LastIndex(fname, "."); i != -1 && i < len(fname)-1 {
			tname = fname[:i+1] + c.Field
		}
		f, ok := r.fields[tname]
		if !ok {
			return nil, fmt.Errorf("unknown target %q", tname)
		}
		return &ConditionRequireIf{Target: f, ValueStr: c.Value}, nil

	case tagCondEnum:
		return ConditionEnum(c), nil

	default:
		return nil, fmt.Errorf("unknown tag condition type %T", cond)
	}
}

func combineConditions(parent *ConditionRequireIf, current interface{}) interface{} {
	if parent == nil {
		return current
	}

	switch c := current.(type) {
	case ConditionRequired:
		if !c { // Required(false)
			return c
		}
		return parent

	case *ConditionRequireIf:
		return c

	case ConditionEnum:
		return &ConditionRequireIfCombined{
			First:  parent,
			Second: c,
		}
	}

	panic(fmt.Errorf("unknown condition %T", current))
}

func findUnmarshalerNameForField(field *reflect.StructField) (unmarshal.UnmarshalerName, bool) {
	const maxPartsCount = 2
	if explicit := field.Tag.Get(unmarshalerTagName); explicit != "" {
		parts := strings.SplitN(explicit, " ", maxPartsCount)
		if len(parts) == maxPartsCount {
			return unmarshal.UnmarshalerName{Package: parts[0], Type: parts[1]}, true
		}
		return unmarshal.UnmarshalerName{Package: field.Type.PkgPath(), Type: parts[0]}, true
	}

	return registry.FindUnmarshalerNameForType(field.Type)
}

func parseTypeName(t reflect.Type) (string, string) {
	name := t.String()
	if i := strings.LastIndex(name, "."); i != -1 {
		return name[:i], name[i+1:]
	}
	panic(fmt.Errorf("strange type name: %q", name))
}

func dereferenceT(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		return t.Elem()
	}
	return t
}

func dereferenceV(v *reflect.Value) *reflect.Value {
	if v.Kind() == reflect.Ptr {
		v := v.Elem()
		return &v
	}
	return v
}

func newValueOfType(t reflect.Type) reflect.Value {
	var v reflect.Value
	if t.Kind() == reflect.Ptr {
		v = reflect.New(t.Elem())
	} else {
		v = reflect.Zero(t)
	}
	return v
}
