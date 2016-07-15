package storeutils

import (
	"crypto/sha1"
	"encoding/hex"
	"strings"
)

func maybeSha1(hash bool, key string) string {
	if hash {
		key = Sha1(key)
	}
	return key
}

// Sha1 computes the sha1 hash of the input.
func Sha1(k string) string {
	hk := sha1.Sum([]byte(k))
	return hex.EncodeToString(hk[:])
}

// Key computes a full key from the input namespace.  The last element of
// the namespace represents a potentially very long unique identifier, and
// it is optionally Sha1 hashed if the hash flag is set.
//
// For example, Key(true, "n0", "n1", "id") -> n0:n1:Sha1("id")
func Key(hash bool, namespace ...string) string {
	n := len(namespace)
	if n == 0 {
		return ""
	}
	id := namespace[n-1]
	namespace[n-1] = maybeSha1(hash, id)
	return strings.Join(namespace, ":")
}
