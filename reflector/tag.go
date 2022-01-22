package reflector

import (
	"regexp"
	"strings"
)

var (
	tagRx      = regexp.MustCompile(`(?s)^(\w+\*?|\*)(?:,(\w+)(?:\((.+)\))?)?$`)
	optReqIfRx = regexp.MustCompile(`(?s)^(\w+)=(.+)$`)
)

type tagInfo struct {
	EnvName string
	Cond    interface{}
}

type tagCondOptional struct{}

type tagCondIf struct {
	Field string
	Value string
}

type tagCondEnum []string

func parseEnvTag(str string) (*tagInfo, error) {
	matches := tagRx.FindStringSubmatch(str)
	if matches == nil {
		return nil, ErrInvalidEnvTag
	}

	tag := &tagInfo{
		EnvName: matches[1],
	}

	if opt := matches[2]; opt != "" {
		val := matches[3]

		switch opt {
		case "optional":
			if val != "" {
				return nil, ErrInvalidEnvTag
			}
			tag.Cond = tagCondOptional{}

		case "reqif":
			m := optReqIfRx.FindStringSubmatch(val)
			if m == nil {
				return nil, ErrInvalidEnvTag
			}
			tag.Cond = &tagCondIf{
				Field: m[1],
				Value: m[2],
			}

		case "enum":
			enum := strings.Split(val, ",")
			if len(enum) == 0 {
				return nil, ErrInvalidEnvTag
			}
			for _, v := range enum {
				if v == "" {
					return nil, ErrInvalidEnvTag
				}
			}
			tag.Cond = tagCondEnum(enum)
		default:
			return nil, ErrInvalidEnvTag
		}
	}

	return tag, nil
}
