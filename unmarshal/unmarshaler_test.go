package unmarshal_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCase struct {
	Name      string
	Input     string
	Result    interface{}
	AssertErr assertErrorFunc
}

type assertErrorFunc func(t *testing.T, err error)

func assertEqualError(str string) assertErrorFunc {
	return func(t *testing.T, err error) {
		t.Helper()
		assert.EqualError(t, err, str)
	}
}
