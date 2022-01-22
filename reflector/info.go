package reflector

import (
	"github.com/albenik/huenv/unmarshal"
)

type TargetInfo struct {
	Target    *Target
	Condition interface{}
}

type Target struct {
	Field       string
	Unmarshaler unmarshal.UnmarshalerName
}

type ConditionRequired bool

type ConditionRequireIf struct {
	Target   *Target
	ValueStr string
}

type ConditionRequireIfCombined struct {
	First  *ConditionRequireIf
	Second interface{}
}

type ConditionEnum []string
