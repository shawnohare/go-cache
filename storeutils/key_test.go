package storeutils_test

import (
	"testing"

	"github.com/shawnohare/go-store/storeutils"
	"github.com/stretchr/testify/assert"
)

func TestKey(t *testing.T) {
	tests := []struct {
		ns   []string
		hash bool
		out  string
	}{
		{
			nil,
			true,
			"",
		},
		{
			nil,
			false,
			"",
		},
		{
			[]string{"key"},
			false,
			"key",
		},
		{
			[]string{"key"},
			true,
			storeutils.Sha1("key"),
		},
		{
			[]string{"n1", "n2"},
			false,
			"n1:n2",
		},
		{
			[]string{"n1", "n2"},
			true,
			"n1:" + storeutils.Sha1("n2"),
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.out, storeutils.Key(tt.hash, tt.ns...))
	}
}
