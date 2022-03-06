package unmarshal_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/albenik/huenv/unmarshal"
)

func TestStringsSlice_Unmarshal(t *testing.T) {
	testCases := []*testCase{{
		Name:   "Single",
		Input:  "foo",
		Result: []string{"foo"},
	}, {
		Name:   "Multi",
		Input:  "foo,bar",
		Result: []string{"foo", "bar"},
	}}

	for _, c := range testCases {
		c := c
		t.Run(c.Name, func(t *testing.T) {
			var val []string
			u := new(unmarshal.StringsSlice)
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
