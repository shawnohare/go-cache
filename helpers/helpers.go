package helpers

import (
	"crypto/sha1"
	"encoding/hex"
	"strings"
)

// Sha1 computes the sha1 hash of the input key.
func Sha1(key string) string {
	hk := sha1.Sum([]byte(key))
	return hex.EncodeToString(hk[:])
}

func maybeSha1(key string, hash bool) string {
	if hash {
		key = Sha1(key)
	}
	return key
}

// Key computes a full key from the input and namespace.  If the hash
// flag is set, the final component k is sha1 hashed.
//
// For example, Key(true, "key", "n0", "n1"}) -> n0:n1:Sha1("key")
func Key(hash bool, k string, namespace ...string) string {
	namespace = append(namespace, maybeSha1(k, hash))
	return strings.Join(namespace, ":")
}
