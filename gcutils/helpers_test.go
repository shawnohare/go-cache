package gcutils_test

import (
	"testing"

	"github.com/shawnohare/go-cache/gcutils"
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
			gcutils.Sha1(""),
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
			gcutils.Sha1("key"),
		},
		{
			"key",
			[]string{},
			true,
			gcutils.Sha1("key"),
		},
		{
			"key",
			[]string{"n0", "n1"},
			true,
			"n0:n1:" + gcutils.Sha1("key"),
		},
		{
			"key",
			[]string{"n0", "n1"},
			false,
			"n0:n1:key",
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.out, gcutils.Key(tt.hash, tt.ns, tt.k))
	}
}
