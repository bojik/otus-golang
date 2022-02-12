package hw09structvalidator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseTag(t *testing.T) {
	tcs := []struct {
		tag string
		exp []TagData
	}{
		{"len:36", []TagData{{"len", []string{"36"}, "len:36"}}},
		{
			"min:18|max:50",
			[]TagData{
				{"min", []string{"18"}, "min:18"},
				{"max", []string{"50"}, "max:50"},
			},
		},
		{
			"in:admin,stuff",
			[]TagData{{"in", []string{"admin", "stuff"}, "in:admin,stuff"}},
		},
		{"required", []TagData{{"required", nil, "required"}}},
		{
			"regexp:^\\w+@\\w+\\.\\w+$",
			[]TagData{
				{"regexp", []string{"^\\w+@\\w+\\.\\w+$"}, "regexp:^\\w+@\\w+\\.\\w+$"},
			},
		},
		{"", nil},
	}
	for _, tc := range tcs {
		tc := tc
		t.Run("parse "+tc.tag, func(t *testing.T) {
			require.Equal(t, tc.exp, ParseTag(tc.tag))
		})
	}
}
