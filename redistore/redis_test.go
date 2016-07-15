package redistore_test

import (
	"encoding/json"
	"flag"
	"testing"

	"github.com/garyburd/redigo/redis"
	"github.com/shawnohare/go-store/redistore"

	. "gopkg.in/check.v1"
)

var local = flag.Bool("local", false, "Include tests that connect to local redis.")

// Test object to store in redis.
type testObj struct {
	X int
}

func Test(t *testing.T) { TestingT(t) }

// RedisSuite performs tests that involve actually testing to a Redis instance
// on the default port.
type RedisSuite struct {
	store *redistore.Store
}

var _ = Suite(new(RedisSuite))

func (s *RedisSuite) SetUpSuite(c *C) {
	if !*local {
		c.Skip("-local not provided")
	}
	// Connect to local redis.

	st := &redistore.Store{Pool: redistore.NewPool(), HashKeys: true}
	s.store = st
	_ = st.Set(st.Key("test", "key"), "val")
	_ = st.Set(st.Key("test", "intkey"), 1)
	_ = st.Set(st.Key("test", "objkey"), testObj{X: 2})
	_ = st.HSet(st.Key("test", "hkey"), "field", "hval")
}

func (s *RedisSuite) SuiteTearDown(c *C) {
	_ = s.store.Del(s.store.Key("test", "key"))
	_ = s.store.Del(s.store.Key("test", "intkey"))
	_ = s.store.Del(s.store.Key("test", "objkey"))
	_ = s.store.Del(s.store.Key("test", "hkey"))
	_ = s.store.HDel(s.store.Key("test", "hkey"), "field")
}

func (s *RedisSuite) TeststoreGetString(c *C) {
	v, ok, err := s.store.Get(s.store.Key("test", "key"))
	c.Assert(err, IsNil)
	c.Assert(ok, Equals, true)
	c.Assert(string(v), Equals, `val`)
}

func (s *RedisSuite) TeststoreGetInt(c *C) {
	v, ok, err := s.store.Get(s.store.Key("test", "intkey"))
	c.Assert(err, IsNil)
	c.Assert(ok, Equals, true)
	vint, _ := redis.Int(v, err)
	c.Assert(vint, Equals, 1)
}

func (s *RedisSuite) TeststoreGetObj(c *C) {
	v, ok, err := s.store.Get(s.store.Key("test", "objkey"))
	c.Assert(err, IsNil)
	c.Assert(ok, Equals, true)

	expected := testObj{X: 2}
	var actual testObj
	_ = json.Unmarshal(v, &actual)
	c.Assert(actual, Equals, expected)
}

func (s *RedisSuite) TeststoreGetNoExist(c *C) {
	v, ok, err := s.store.Get("redistore-key-does-not-exist")
	c.Assert(err, IsNil)
	c.Assert(ok, Equals, false)
	c.Assert(v, IsNil)
}

func (s *RedisSuite) TeststoreHGetString(c *C) {
	v, ok, err := s.store.HGet(s.store.Key("test", "hkey"), "field")
	c.Assert(err, IsNil)
	c.Assert(ok, Equals, true)
	c.Assert(string(v), Equals, `hval`)
}

func (s *RedisSuite) TeststoreHGetNoExist(c *C) {
	v, ok, err := s.store.HGet("redistore-key-does-not-exist", "field")
	c.Assert(err, IsNil)
	c.Assert(ok, Equals, false)
	c.Assert(v, IsNil)
}

func (s *RedisSuite) TeststoreSetObj(c *C) {
	err := s.store.Set("new obj key", testObj{X: 3})
	defer func() { _ = s.store.Del("new obj key") }()
	c.Assert(err, IsNil)
}

func (s *RedisSuite) TeststoreDel(c *C) {
	_ = s.store.Set("new key", 1)
	err := s.store.Del("new key")
	c.Assert(err, IsNil)
	_, ok, err := s.store.Get("new key")
	c.Assert(err, IsNil)
	c.Assert(ok, Equals, false)
}

func (s *RedisSuite) TeststoreHDel(c *C) {
	k := "new hkey"
	f := "field"
	_ = s.store.HSet(k, f, 1)
	err := s.store.HDel(k, f)
	c.Assert(err, IsNil)
	_, ok, err := s.store.HGet(k, f)
	c.Assert(ok, Equals, false)
}
