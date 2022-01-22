package generator_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/albenik/huenv/generator"
)

func TestGenerateProgram(t *testing.T) {
	src, err := generator.GenerateProgram("example.net/lib/config", "Config")
	require.NoError(t, err)

	formatted := bytes.NewBuffer(nil)

	err = generator.FormatSource(formatted, bytes.NewReader(src))
	if !assert.NoError(t, err) {
		var de generator.DetailedError
		if errors.As(err, &de) {
			t.Log(de.Details())
		}
		t.FailNow()
	}
	assert.Equal(t, formatted.String(), string(src))
}
