package unmarshal_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/albenik/huenv/unmarshal"
)

func TestTime_Unmarshal(t *testing.T) {
	assertTimeParseError := func(t *testing.T, err error) {
		t.Helper()
		aserr := new(time.ParseError)
		assert.ErrorAs(t, err, &aserr)
	}

	testCases := []*testCase{{
		Name:   "Valid",
		Input:  time.Date(2022, time.January, 1, 12, 30, 0, 0, time.UTC).Format(time.RFC3339),
		Result: time.Date(2022, time.January, 1, 12, 30, 0, 0, time.UTC),
	}, {
		Name:      "Invalid",
		Input:     "invalid",
		AssertErr: assertTimeParseError,
	}}

	for _, c := range testCases {
		c := c
		t.Run(c.Name, func(t *testing.T) {
			var val time.Time
			u := new(unmarshal.Time)
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

func TestDuration_Unmarshal(t *testing.T) {
	testCases := []*testCase{{
		Name:   "5s",
		Input:  "5s",
		Result: 5 * time.Second,
	}, {
		Name:   "1m",
		Input:  "1m",
		Result: time.Minute,
	}, {
		Name:      "Invalid(Text)",
		Input:     "foo",
		AssertErr: assertEqualError("time: invalid duration \"foo\""),
	}, {
		Name:      "Invalid(Num)",
		Input:     "5",
		AssertErr: assertEqualError("time: missing unit in duration \"5\""),
	}}

	for _, c := range testCases {
		c := c
		t.Run(c.Name, func(t *testing.T) {
			var val time.Duration
			u := new(unmarshal.Duration)
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
