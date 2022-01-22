package unmarshal_test

import (
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/albenik/huenv/unmarshal"
)

func TestInt_Unmarshal(t *testing.T) {
	assertNumError := func(t *testing.T, err error) {
		t.Helper()
		aserr := new(strconv.NumError)
		assert.ErrorAs(t, err, &aserr)
	}

	testCases := []*testCase{{
		Name:   "Zero",
		Input:  "0",
		Result: 0,
	}, {
		Name:   "Valid/Max",
		Input:  strconv.FormatInt(math.MaxInt, 10),
		Result: math.MaxInt,
	}, {
		Name:      "Invalid/Str",
		Input:     "invalid",
		AssertErr: assertNumError,
	}}

	for _, c := range testCases {
		c := c
		t.Run(c.Name, func(t *testing.T) {
			var val int
			u := new(unmarshal.Int)
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

func TestInt8_Unmarshal(t *testing.T) {
	assertNumError := func(t *testing.T, err error) {
		t.Helper()
		aserr := new(strconv.NumError)
		assert.ErrorAs(t, err, &aserr)
	}

	testCases := []*testCase{{
		Name:   "Zero",
		Input:  "0",
		Result: int8(0),
	}, {
		Name:   "Valid/Max",
		Input:  strconv.FormatInt(math.MaxInt8, 10),
		Result: int8(math.MaxInt8),
	}, {
		Name:      "Invalid/Str",
		Input:     "invalid",
		AssertErr: assertNumError,
	}}

	for _, c := range testCases {
		c := c
		t.Run(c.Name, func(t *testing.T) {
			var val int8
			u := new(unmarshal.Int8)
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

func TestInt16_Unmarshal(t *testing.T) {
	assertNumError := func(t *testing.T, err error) {
		t.Helper()
		aserr := new(strconv.NumError)
		assert.ErrorAs(t, err, &aserr)
	}

	testCases := []*testCase{{
		Name:   "Zero",
		Input:  "0",
		Result: int16(0),
	}, {
		Name:   "Valid/Max",
		Input:  strconv.FormatInt(math.MaxInt16, 10),
		Result: int16(math.MaxInt16),
	}, {
		Name:      "Invalid/Str",
		Input:     "invalid",
		AssertErr: assertNumError,
	}}

	for _, c := range testCases {
		c := c
		t.Run(c.Name, func(t *testing.T) {
			var val int16
			u := new(unmarshal.Int16)
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

func TestInt32_Unmarshal(t *testing.T) {
	assertNumError := func(t *testing.T, err error) {
		t.Helper()
		aserr := new(strconv.NumError)
		assert.ErrorAs(t, err, &aserr)
	}

	testCases := []*testCase{{
		Name:   "Zero",
		Input:  "0",
		Result: int32(0),
	}, {
		Name:   "Valid/Max",
		Input:  strconv.FormatInt(math.MaxInt32, 10),
		Result: int32(math.MaxInt32),
	}, {
		Name:      "Invalid/Str",
		Input:     "invalid",
		AssertErr: assertNumError,
	}}

	for _, c := range testCases {
		c := c
		t.Run(c.Name, func(t *testing.T) {
			var val int32
			u := new(unmarshal.Int32)
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

func TestInt64_Unmarshal(t *testing.T) {
	assertNumError := func(t *testing.T, err error) {
		t.Helper()
		aserr := new(strconv.NumError)
		assert.ErrorAs(t, err, &aserr)
	}

	testCases := []*testCase{{
		Name:   "Zero",
		Input:  "0",
		Result: int64(0),
	}, {
		Name:   "Valid/Max",
		Input:  strconv.FormatInt(math.MaxInt64, 10),
		Result: int64(math.MaxInt64),
	}, {
		Name:      "Invalid/Str",
		Input:     "invalid",
		AssertErr: assertNumError,
	}}

	for _, c := range testCases {
		c := c
		t.Run(c.Name, func(t *testing.T) {
			var val int64
			u := new(unmarshal.Int64)
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

func TestUint_Unmarshal(t *testing.T) {
	assertNumError := func(t *testing.T, err error) {
		t.Helper()
		aserr := new(strconv.NumError)
		assert.ErrorAs(t, err, &aserr)
	}

	testCases := []*testCase{{
		Name:   "Zero",
		Input:  "0",
		Result: uint(0),
	}, {
		Name:   "Valid/Max",
		Input:  strconv.FormatUint(math.MaxUint, 10),
		Result: uint(math.MaxUint),
	}, {
		Name:      "Invalid/Str",
		Input:     "invalid",
		AssertErr: assertNumError,
	}}

	for _, c := range testCases {
		c := c
		t.Run(c.Name, func(t *testing.T) {
			var val uint
			u := new(unmarshal.Uint)
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

func TestUint8_Unmarshal(t *testing.T) {
	assertNumError := func(t *testing.T, err error) {
		t.Helper()
		aserr := new(strconv.NumError)
		assert.ErrorAs(t, err, &aserr)
	}

	testCases := []*testCase{{
		Name:   "Zero",
		Input:  "0",
		Result: uint8(0),
	}, {
		Name:   "Valid/Max",
		Input:  strconv.FormatUint(math.MaxUint8, 10),
		Result: uint8(math.MaxUint8),
	}, {
		Name:      "Invalid/Str",
		Input:     "invalid",
		AssertErr: assertNumError,
	}}

	for _, c := range testCases {
		c := c
		t.Run(c.Name, func(t *testing.T) {
			var val uint8
			u := new(unmarshal.Uint8)
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

func TestUint16_Unmarshal(t *testing.T) {
	assertNumError := func(t *testing.T, err error) {
		t.Helper()
		aserr := new(strconv.NumError)
		assert.ErrorAs(t, err, &aserr)
	}

	testCases := []*testCase{{
		Name:   "Zero",
		Input:  "0",
		Result: uint16(0),
	}, {
		Name:   "Valid/Max",
		Input:  strconv.FormatUint(math.MaxUint16, 10),
		Result: uint16(math.MaxUint16),
	}, {
		Name:      "Invalid/Str",
		Input:     "invalid",
		AssertErr: assertNumError,
	}}

	for _, c := range testCases {
		c := c
		t.Run(c.Name, func(t *testing.T) {
			var val uint16
			u := new(unmarshal.Uint16)
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

func TestUint32_Unmarshal(t *testing.T) {
	assertNumError := func(t *testing.T, err error) {
		t.Helper()
		aserr := new(strconv.NumError)
		assert.ErrorAs(t, err, &aserr)
	}

	testCases := []*testCase{{
		Name:   "Zero",
		Input:  "0",
		Result: uint32(0),
	}, {
		Name:   "Valid/Max",
		Input:  strconv.FormatUint(math.MaxUint32, 10),
		Result: uint32(math.MaxUint32),
	}, {
		Name:      "Invalid/Str",
		Input:     "invalid",
		AssertErr: assertNumError,
	}}

	for _, c := range testCases {
		c := c
		t.Run(c.Name, func(t *testing.T) {
			var val uint32
			u := new(unmarshal.Uint32)
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

func TestUint64_Unmarshal(t *testing.T) {
	assertNumError := func(t *testing.T, err error) {
		t.Helper()
		aserr := new(strconv.NumError)
		assert.ErrorAs(t, err, &aserr)
	}

	testCases := []*testCase{{
		Name:   "Zero",
		Input:  "0",
		Result: uint64(0),
	}, {
		Name:   "Valid/Max",
		Input:  strconv.FormatUint(math.MaxUint64, 10),
		Result: uint64(math.MaxUint64),
	}, {
		Name:      "Invalid/Str",
		Input:     "invalid",
		AssertErr: assertNumError,
	}}

	for _, c := range testCases {
		c := c
		t.Run(c.Name, func(t *testing.T) {
			var val uint64
			u := new(unmarshal.Uint64)
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
