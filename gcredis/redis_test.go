package gcredis_test

import (
	"encoding/json"
	"flag"
	"testing"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/shawnohare/go-cache/gcredis"
	"github.com/shawnohare/go-cache/gcutils"

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
	_ = s.cache.Set(gcutils.Namespace("test"), "key", "val")
	_ = s.cache.Set(gcutils.Namespace("test"), "intkey", 1)
	_ = s.cache.Set(gcutils.Namespace("test"), "objkey", testObj{X: 2})
	_ = s.cache.HSet(gcutils.Namespace("test"), "hkey", "field", "hval")
}

func (s *RedisSuite) SuiteTearDown(c *C) {
	_ = s.cache.Del([]string{"test"}, "key")
	_ = s.cache.Del([]string{"test"}, "intkey")
	_ = s.cache.Del([]string{"test"}, "objkey")
	_ = s.cache.HDel([]string{"test"}, "hkey", "field")
	_ = s.cache.Del([]string{"test"}, "hkey")
}

func (s *RedisSuite) TestCacheGetString(c *C) {
	v, ok, err := s.cache.Get([]string{"test"}, "key")
	c.Assert(err, IsNil)
	c.Assert(ok, Equals, true)
	c.Assert(string(v), Equals, `val`)
}

func (s *RedisSuite) TestCacheGetInt(c *C) {
	v, ok, err := s.cache.Get([]string{"test"}, "intkey")
	c.Assert(err, IsNil)
	c.Assert(ok, Equals, true)
	vint, _ := redis.Int(v, err)
	c.Assert(vint, Equals, 1)
}

func (s *RedisSuite) TestCacheGetObj(c *C) {
	v, ok, err := s.cache.Get([]string{"test"}, "objkey")
	c.Assert(err, IsNil)
	c.Assert(ok, Equals, true)

	expected := testObj{X: 2}
	var actual testObj
	_ = json.Unmarshal(v, &actual)
	c.Assert(actual, Equals, expected)
}

func (s *RedisSuite) TestCacheGetNoExist(c *C) {
	v, ok, err := s.cache.Get([]string{"test"}, "gcredis-key-does-not-exist")
	c.Assert(err, IsNil)
	c.Assert(ok, Equals, false)
	c.Assert(v, IsNil)
}

func (s *RedisSuite) TestCacheHGetString(c *C) {
	v, ok, err := s.cache.HGet([]string{"test"}, "hkey", "field")
	c.Assert(err, IsNil)
	c.Assert(ok, Equals, true)
	c.Assert(string(v), Equals, `hval`)
}

func (s *RedisSuite) TestCacheHGetNoExist(c *C) {
	v, ok, err := s.cache.HGet([]string{"test"}, "gcredis-key-does-not-exist", "field")
	c.Assert(err, IsNil)
	c.Assert(ok, Equals, false)
	c.Assert(v, IsNil)
}

func (s *RedisSuite) TestCacheSetObj(c *C) {
	err := s.cache.Set(nil, "new obj key", testObj{X: 3})
	defer func() { _ = s.cache.Del(nil, "new obj key") }()
	c.Assert(err, IsNil)
}

func (s *RedisSuite) TestCacheDel(c *C) {
	_ = s.cache.Set(nil, "new key", 1)
	err := s.cache.Del(nil, "new key")
	c.Assert(err, IsNil)
	_, ok, err := s.cache.Get(nil, "new key")
	c.Assert(err, IsNil)
	c.Assert(ok, Equals, false)
}

func (s *RedisSuite) TestCacheHDel(c *C) {
	k := "new hkey"
	f := "field"
	ns := []string{"test"}
	_ = s.cache.HSet(ns, k, f, 1)
	err := s.cache.HDel(ns, k, f)
	c.Assert(err, IsNil)
	_, ok, err := s.cache.HGet(ns, k, f)
	c.Assert(ok, Equals, false)
}
