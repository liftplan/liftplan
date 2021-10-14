package fto

import (
	"net/url"
	"testing"
)

func TestIsChecked(t *testing.T) {
	t.Parallel()
	v := url.Values{
		"foo": []string{"true"},
		"bar": []string{"false"},
		"bam": []string{},
	}
	tt := []struct {
		key      string
		vals     url.Values
		expected bool
	}{
		{"foo", v, true},
		{"bar", v, false},
		{"baz", v, false},
		{"bam", v, false},
	}

	for _, test := range tt {
		k, v, e := test.key, test.vals, test.expected
		if isChecked(k, v) != e {
			t.Errorf("key: %v, value: %v, expected: %v", k, v, e)
		}
	}
}
