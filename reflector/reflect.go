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

type Result struct {
	ConfigPkg  string
	ConfigType string
	Packages   map[string]struct{}
	Envs       map[string]*TargetInfo
}

type Reflector struct {
	pkgs   map[string]struct{}
	envs   map[string]*TargetInfo
	fields map[string]*Target
}

func New() *Reflector {
	return &Reflector{
		pkgs:   make(map[string]struct{}),
		envs:   make(map[string]*TargetInfo),
		fields: make(map[string]*Target),
	}
}

func (r *Reflector) Reflect(conf interface{}) (*Result, error) {
	t := reflect.TypeOf(conf)
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		return nil, InvalidConfigTypeError(t.String())
	}

	ts := t.Elem()

	if err := r.processStruct(ts, "", "", nil); err != nil {
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

func (r *Reflector) processStruct(typ reflect.Type, envPrefix, fieldPrefix string, parentCond *ConditionRequireIf) error {
	typ.String()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		longFieldName := field.Name
		if fieldPrefix != "" {
			longFieldName = fmt.Sprintf("%s.%s", fieldPrefix, longFieldName)
		}

		uname, hasUnmarshaler := findUnmarshalerNameForField(&field)

		if !hasUnmarshaler {
			if processed, err := r.tryProcessSublevel(&field, longFieldName, envPrefix); err != nil {
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

		info := &TargetInfo{
			Target: &Target{
				Field:       longFieldName,
				Unmarshaler: uname,
			},
			Condition: combineConditions(parentCond, cond),
		}

		r.pkgs[uname.Package] = struct{}{}
		r.envs[envName] = info
		r.fields[longFieldName] = info.Target
	}

	return nil
}

func (r *Reflector) tryProcessSublevel(field *reflect.StructField, longFieldName, envPrefix string) (bool, error) {
	t := dereference(field.Type)
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

	if err := r.processStruct(t, envPrefix, longFieldName, condition); err != nil {
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

func findUnmarshalerNameForField(field *reflect.StructField) (unmarshal.UnmarshalerName, bool) {
	if explicit := field.Tag.Get(unmarshalerTagName); explicit != "" {
		parts := strings.SplitN(explicit, " ", 2)
		if len(parts) == 2 {
			return unmarshal.UnmarshalerName{Package: parts[0], Type: parts[1]}, true
		}
		return unmarshal.UnmarshalerName{Package: field.Type.PkgPath(), Type: parts[0]}, true
	}

	return registry.FindUnmarshalerNameForType(field.Type)
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

func dereference(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		return t.Elem()
	}
	return t
}

func parseTypeName(t reflect.Type) (string, string) {
	name := t.String()
	if i := strings.LastIndex(name, "."); i != -1 {
		return name[:i], name[i+1:]
	}
	panic(fmt.Errorf("strange type name: %q", name))
}
