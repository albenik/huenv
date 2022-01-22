package unmarshal_test

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/albenik/huenv/unmarshal"
)

func TestURL_Unmarshal(t *testing.T) {
	testCases := []*testCase{{
		Name:  "Valid",
		Input: "https://example.com:8888/foo/bar",
		Result: &url.URL{
			Scheme:      "https",
			Opaque:      "",
			User:        nil,
			Host:        "example.com:8888",
			Path:        "/foo/bar",
			RawPath:     "",
			ForceQuery:  false,
			RawQuery:    "",
			Fragment:    "",
			RawFragment: "",
		},
	}, {
		Name:      "Invalid",
		Input:     "://foo",
		AssertErr: assertEqualError("parse \"://foo\": missing protocol scheme"),
	}}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			var val *url.URL
			u := new(unmarshal.URL)
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
