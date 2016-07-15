package storeutils_test

import (
	"testing"

	"github.com/shawnohare/go-store/storeutils"
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
			storeutils.Sha1(""),
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
			storeutils.Sha1("key"),
		},
		{
			"key",
			[]string{},
			true,
			storeutils.Sha1("key"),
		},
		{
			"key",
			[]string{"n0", "n1"},
			true,
			"n0:n1:" + storeutils.Sha1("key"),
		},
		{
			"key",
			[]string{"n0", "n1"},
			false,
			"n0:n1:key",
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.out, storeutils.Key(tt.hash, tt.ns, tt.k))
	}
}
