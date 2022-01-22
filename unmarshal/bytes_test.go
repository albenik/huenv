package unmarshal_test

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/albenik/huenv/unmarshal"
)

func TestBytes_Unmarshal(t *testing.T) {
	assertCorruptInputError := func(t *testing.T, err error) {
		t.Helper()
		var aserr base64.CorruptInputError
		assert.ErrorAs(t, err, &aserr)
	}

	testCases := []*testCase{{
		Name:   "Valid",
		Input:  base64.StdEncoding.EncodeToString([]byte{1, 2, 3, 4, 5}),
		Result: []byte{1, 2, 3, 4, 5},
	}, {
		Name:      "Invalid",
		Input:     "ЪЫФ",
		AssertErr: assertCorruptInputError,
	}}

	for _, c := range testCases {
		c := c
		t.Run(c.Name, func(t *testing.T) {
			var val []byte

			u := new(unmarshal.Bytes)
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
