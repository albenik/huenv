package unmarshal_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/albenik/huenv/unmarshal"
)

func TestBool_Unmarshal(t *testing.T) {
	assertNumError := func(t *testing.T, err error) {
		t.Helper()
		aserr := new(strconv.NumError)
		assert.ErrorAs(t, err, &aserr)
	}

	testCases := []*testCase{{
		Name:   "False",
		Input:  "false",
		Result: false,
	}, {
		Name:   "True",
		Input:  "true",
		Result: true,
	}, {
		Name:      "Invalid",
		Input:     "foo",
		AssertErr: assertNumError,
	}}

	for _, c := range testCases {
		c := c
		t.Run(c.Name, func(t *testing.T) {
			var val bool
			u := new(unmarshal.Bool)
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
