package gcredis

import (
	"encoding/json"

	"github.com/garyburd/redigo/redis"
	"github.com/shawnohare/go-cache/helpers"
)

// type Pool interface {
// 	Get() redis.Conn
// }

// Cache is a thin cache wrapper over a redigo.Pool that adds some basic typing
// to the redigo commands.
//
// If HashKeys is true, the Cache will sha1 hash the final component of
// all namespaced keys.  For example, the key, namespace pair
// (k, []string{n0, n1}) will map to n0:n1:k -> n0:n1:sha1(k).
// Otherwise, it maps to n0:n1:k.
type Cache struct {
	Pool     *redis.Pool
	HashKeys bool
}

// wrapper for closing a connection that ignores errors.  This is defined
// primarily to avoid handling the error.
func (c *Cache) Close(conn redis.Conn) {
	_ = conn.Close()
}

// Marshal wraps the redis.Bytes function to convert a value to cache
// into a byte slice.
func (c *Cache) Marshal(v interface{}) ([]byte, error) {
	var bs []byte
	var err error
	switch t := v.(type) {
	case string:
		bs = []byte(t)
	case []byte:
		bs = t
	default:
		bs, err = json.Marshal(t)
	}
	return bs, err
}

func (c *Cache) Unmarshal(response interface{}, err error) ([]byte, bool, error) {
	if response == nil && err == nil { // value does not exist
		return nil, false, nil
	}
	bs, err := redis.Bytes(response, err)
	if err != nil {
		return nil, false, err
	}

	return bs, true, nil
}

// Set saves the (key, value) pair in Cache using the key
// helpers.Key(k, namespace...).
func (c *Cache) Set(k string, value interface{}, namespace ...string) error {
	conn := c.Pool.Get()
	defer c.Close(conn)

	data, err := c.Marshal(value)
	if err != nil {
		return err
	}
	_, err = conn.Do("SET", helpers.Key(c.HashKeys, k, namespace...), data)
	return err
}

// HSet stores the (field, value) pair in the hash keyed by helpers.Key(k, namespace).
func (c *Cache) HSet(k string, field string, value interface{}, namespace ...string) error {
	conn := c.Pool.Get()
	defer c.Close(conn)

	data, err := c.Marshal(value)
	if err != nil {
		return err
	}
	_, err = conn.Do("HMSET", helpers.Key(c.HashKeys, k, namespace...), field, data)
	return err
}

// Get the value stored at helpers.Key(k, namespace).
func (c *Cache) Get(k string, namespace ...string) ([]byte, bool, error) {
	conn := c.Pool.Get()
	defer c.Close(conn)
	return c.Unmarshal(conn.Do("GET", helpers.Key(c.HashKeys, k, namespace...)))
}

// HGet wraps the Cache Hash Get function.
func (c *Cache) HGet(k string, field string, namespace ...string) ([]byte, bool, error) {
	conn := c.Pool.Get()
	defer c.Close(conn)
	return c.Unmarshal(conn.Do("HGET", helpers.Key(c.HashKeys, k, namespace...), field))
}

// Del deletes the value stored at helpers.Key(k, namespace)
func (c *Cache) Del(k string, namespace ...string) error {
	conn := c.Pool.Get()
	defer c.Close(conn)
	_, err := conn.Do("DEL", helpers.Key(c.HashKeys, k, namespace...))
	return err
}

// HDel deletes the field from the hash keyed by helpers.Key(k, namespace)
func (c *Cache) HDel(k string, field string, namespace ...string) error {
	conn := c.Pool.Get()
	defer c.Close(conn)
	_, err := conn.Do("HDEL", helpers.Key(c.HashKeys, k, namespace...), field)
	return err
}
