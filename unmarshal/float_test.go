package unmarshal_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/albenik/huenv/unmarshal"
)

func TestFloat32_Unmarshal(t *testing.T) {
	assertNumError := func(t *testing.T, err error) {
		t.Helper()
		aserr := new(strconv.NumError)
		assert.ErrorAs(t, err, &aserr)
	}

	testCases := []*testCase{{
		Name:   "Valid/Float",
		Input:  "1.2345",
		Result: float32(1.2345),
	}, {
		Name:   "Valid/Int",
		Input:  "1",
		Result: float32(1.0),
	}, {
		Name:      "Invalid",
		Input:     "foo",
		AssertErr: assertNumError,
	}}

	for _, c := range testCases {
		c := c
		t.Run(c.Name, func(t *testing.T) {
			var val float32
			u := new(unmarshal.Float32)
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

func TestFloat64_Unmarshal(t *testing.T) {
	assertNumError := func(t *testing.T, err error) {
		t.Helper()
		aserr := new(strconv.NumError)
		assert.ErrorAs(t, err, &aserr)
	}

	testCases := []*testCase{{
		Name:   "Valid/Float",
		Input:  "1.2345",
		Result: 1.2345,
	}, {
		Name:   "Valid/Int",
		Input:  "1",
		Result: 1.0,
	}, {
		Name:      "Invalid",
		Input:     "foo",
		AssertErr: assertNumError,
	}}

	for _, c := range testCases {
		c := c
		t.Run(c.Name, func(t *testing.T) {
			var val float64
			u := new(unmarshal.Float64)
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
