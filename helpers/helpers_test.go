package helpers_test

import (
	"testing"

	"github.com/shawnohare/go-cache/helpers"
	"github.com/stretchr/testify/assert"
)

func TestKey(t *testing.T) {
	tests := []struct {
		k    string
		ns   []string
		hash bool
		out  string
	}{
		{
			"",
			nil,
			true,
			helpers.Sha1(""),
		},
		{
			"",
			nil,
			false,
			"",
		},
		{
			"key",
			nil,
			true,
			helpers.Sha1("key"),
		},
		{
			"key",
			[]string{},
			true,
			helpers.Sha1("key"),
		},
		{
			"key",
			[]string{"n0", "n1"},
			true,
			"n0:n1:" + helpers.Sha1("key"),
		},
		{
			"key",
			[]string{"n0", "n1"},
			false,
			"n0:n1:key",
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.out, helpers.Key(tt.hash, tt.k, tt.ns...))
	}
}
