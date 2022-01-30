package reflector

import (
	"github.com/albenik/huenv/unmarshal"
)

type Result struct {
	ConfigPkg  string
	ConfigType string
	Packages   map[string]struct{}
	Envs       map[string]*Target
}

type Target struct {
	Field     *TargetField
	Condition interface{}
}

type TargetField struct {
	Name        string
	Unmarshaler unmarshal.UnmarshalerName
}

type ConditionRequired bool

type ConditionRequireIf struct {
	Target   *TargetField
	ValueStr string
}

type ConditionRequireIfCombined struct {
	First  *ConditionRequireIf
	Second interface{}
}

type ConditionEnum []string
