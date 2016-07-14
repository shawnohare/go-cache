package gcredis_test

import (
	"encoding/json"
	"flag"
	"testing"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/shawnohare/go-cache/gcredis"

	. "gopkg.in/check.v1"
)

var local = flag.Bool("local", false, "Include tests that connect to local redis.")

// Test object to store in redis.
type testObj struct {
	X int
}

func Test(t *testing.T) { TestingT(t) }

type RedisSuite struct {
	cache *gcredis.Cache
}

var _ = Suite(new(RedisSuite))

func (s *RedisSuite) SetUpSuite(c *C) {
	if !*local {
		c.Skip("-local not provided")
	}
	// Connect to local redis.
	pool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
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

	s.cache = &gcredis.Cache{Pool: pool, HashKeys: true}
	_ = s.cache.Set("key", "val", "test")
	_ = s.cache.Set("intkey", 1, "test")
	_ = s.cache.Set("objkey", testObj{X: 2}, "test")
	_ = s.cache.HSet("hkey", "field", "hval", "test")
}

func (s *RedisSuite) SuiteTearDown(c *C) {
	_ = s.cache.Del("key", "test")
	_ = s.cache.Del("intkey", "test")
	_ = s.cache.Del("objkey", "test")
	_ = s.cache.HDel("hkey", "field", "test")
}

func (s *RedisSuite) TestCacheGetString(c *C) {
	v, ok, err := s.cache.Get("key", "test")
	c.Assert(err, IsNil)
	c.Assert(ok, Equals, true)
	c.Assert(string(v), Equals, `val`)
}

func (s *RedisSuite) TestCacheGetInt(c *C) {
	v, ok, err := s.cache.Get("intkey", "test")
	c.Assert(err, IsNil)
	c.Assert(ok, Equals, true)
	vint, _ := redis.Int(v, err)
	c.Assert(vint, Equals, 1)
}

func (s *RedisSuite) TestCacheGetObj(c *C) {
	v, ok, err := s.cache.Get("objkey", "test")
	c.Assert(err, IsNil)
	c.Assert(ok, Equals, true)

	expected := testObj{X: 2}
	var actual testObj
	_ = json.Unmarshal(v, &actual)
	c.Assert(actual, Equals, expected)
}

func (s *RedisSuite) TestCacheGetNoExist(c *C) {
	v, ok, err := s.cache.Get("gcredis-key-does-not-exist", "test")
	c.Assert(err, IsNil)
	c.Assert(ok, Equals, false)
	c.Assert(v, IsNil)
}

func (s *RedisSuite) TestCacheHGetString(c *C) {
	v, ok, err := s.cache.HGet("hkey", "field", "test")
	c.Assert(err, IsNil)
	c.Assert(ok, Equals, true)
	c.Assert(string(v), Equals, `hval`)
}

func (s *RedisSuite) TestCacheHGetNoExist(c *C) {
	v, ok, err := s.cache.HGet("gcredis-key-does-not-exist", "field", "test")
	c.Assert(err, IsNil)
	c.Assert(ok, Equals, false)
	c.Assert(v, IsNil)
}

func (s *RedisSuite) TestCacheSetObj(c *C) {
	err := s.cache.Set("new obj key", testObj{X: 3})
	defer func() { _ = s.cache.Del("new obj key") }()
	c.Assert(err, IsNil)
}

func (s *RedisSuite) TestCacheDel(c *C) {
	_ = s.cache.Set("new key", 1, "test")
	err := s.cache.Del("new key", "test")
	c.Assert(err, IsNil)
	_, ok, err := s.cache.Get("new key", "test")
	c.Assert(err, IsNil)
	c.Assert(ok, Equals, false)
}

func (s *RedisSuite) TestCacheHDel(c *C) {
	k := "new hkey"
	f := "field"
	n := "test"
	_ = s.cache.HSet(k, f, 1, n)
	err := s.cache.HDel(k, f, n)
	c.Assert(err, IsNil)
	_, ok, err := s.cache.HGet(k, f, n)
	c.Assert(ok, Equals, false)
}
