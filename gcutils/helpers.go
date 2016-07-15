package gcutils

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

func maybeSha1(hash bool, key string) string {
	if hash {
		key = Sha1(key)
	}
	return key
}

// Namespace converts variadic string input into a slice of strings.
func Namespace(ns ...string) []string {
	return ns
}

// Key computes a full key from the input and namespace.  If the hash
// flag is set, the final component k is sha1 hashed.
//
// For example, Key(true, []string{"n0", "n1"}, "key"}) -> n0:n1:Sha1("key")
func Key(hash bool, namespace []string, k string) string {
	namespace = append(namespace, maybeSha1(hash, k))
	return strings.Join(namespace, ":")
}
