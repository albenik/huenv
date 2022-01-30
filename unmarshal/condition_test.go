package unmarshal_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/albenik/huenv/unmarshal"
)

func TestConditionRequireIf_And(t *testing.T) {
	foo := "jaeger"
	err := unmarshal.RequireIf(&foo, "jaeger", new(unmarshal.String)).And(unmarshal.Enum("agent", "collector")).
		Validate("agent")
	require.NoError(t, err)
}
