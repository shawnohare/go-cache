package gcredis

import (
	"encoding/json"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/shawnohare/go-cache/gcutils"
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

// Connect is a helper function that creates a redis.Pool with some standard
// settings. If no url is provided, a Redis instance running at
// localhost:6379 is assumed.  Otherwise, the provided url is used to connect.
func Pool(url ...string) *redis.Pool {
	var actualURL = ":6379"
	if len(url) > 0 {
		actualURL = url[0]
	}

	pool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", actualURL)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	return pool
}

// Key wraps the gcutils.Key function by passing in the cache's HashKeys flag.
func (c *Cache) Key(namespace []string, k string) string {
	return gcutils.Key(c.HashKeys, namespace, k)
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
// gcutils.Key(namespace, k).
func (c *Cache) Set(namespace []string, k string, value interface{}) error {
	conn := c.Pool.Get()
	defer c.Close(conn)

	data, err := c.Marshal(value)
	if err != nil {
		return err
	}
	_, err = conn.Do("SET", c.Key(namespace, k), data)
	return err
}

// HSet stores the (field, value) pair in the hash keyed by gcutils.Key(k, namespace).
func (c *Cache) HSet(namespace []string, k string, field string, value interface{}) error {
	conn := c.Pool.Get()
	defer c.Close(conn)

	data, err := c.Marshal(value)
	if err != nil {
		return err
	}
	_, err = conn.Do("HMSET", c.Key(namespace, k), field, data)
	return err
}

// Get the value stored at gcutils.Key(k, namespace).
func (c *Cache) Get(namespace []string, k string) ([]byte, bool, error) {
	conn := c.Pool.Get()
	defer c.Close(conn)
	return c.Unmarshal(conn.Do("GET", c.Key(namespace, k)))
}

// HGet wraps the Cache Hash Get function.
func (c *Cache) HGet(namespace []string, k string, field string) ([]byte, bool, error) {
	conn := c.Pool.Get()
	defer c.Close(conn)
	return c.Unmarshal(conn.Do("HGET", c.Key(namespace, k), field))
}

// Del deletes the value stored at gcutils.Key(k, namespace)
func (c *Cache) Del(namespace []string, k string) error {
	conn := c.Pool.Get()
	defer c.Close(conn)
	_, err := conn.Do("DEL", c.Key(namespace, k))
	return err
}

// HDel deletes the field from the hash keyed by gcutils.Key(k, namespace)
func (c *Cache) HDel(namespace []string, k string, field string) error {
	conn := c.Pool.Get()
	defer c.Close(conn)
	_, err := conn.Do("HDEL", c.Key(namespace, k), field)
	return err
}
