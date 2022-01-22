package unmarshal_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/albenik/huenv/unmarshal"
)

func TestString_Unmarshal(t *testing.T) {
	testCases := []*testCase{{
		Name:   "Set",
		Input:  "foo",
		Result: "foo",
	}}

	for _, c := range testCases {
		c := c
		t.Run(c.Name, func(t *testing.T) {
			var val string
			u := new(unmarshal.String)
			require.NoError(t, u.SetTarget(&val))

			err := u.Unmarshal(c.Input)
			if c.AssertErr != nil {
				c.AssertErr(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, c.Result, val)
			}
		})
	}
}
