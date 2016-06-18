package main

import (
	"os"

	"github.com/bitly/go-simplejson"
	"github.com/garyburd/redigo/redis"
)

type Cache interface {
	GetBool(key string) bool
	GetString(key string) string
	GetMap(key string) map[string]int
	Set(key string, val interface{})
	Del(key string)
	Close()
}

type CacheFile struct {
	json *simplejson.Json
}

// NewChecker is generator for Checker
func NewCacheFile() *CacheFile {
	var r *os.File
	_, err := os.Stat(".million-timer")
	if err == nil {
		r, _ = os.Open(".million-timer")
	} else {
		r, _ = os.Create(".million-timer")
	}
	json, err := simplejson.NewFromReader(r)
	if err != nil {
		json = simplejson.New()
	}

	return &CacheFile{json: json}
}

func (c *CacheFile) GetBool(key string) bool {
	return c.json.Get(key).MustBool(false)
}

func (c *CacheFile) GetString(key string) string {
	return c.json.Get(key).MustString("")
}

func (c *CacheFile) GetMap(key string) map[string]int {
	m := make(map[string]int)
	for k, _ := range c.json.Get(key).MustMap() {
		m[k] = 1
	}
	return m
}

func (c *CacheFile) Set(key string, val interface{}) {
	c.json.Set(key, val)
}

func (c *CacheFile) Del(key string) {
	c.json.Del(key)
}

func (c *CacheFile) Close() {
	w, _ := os.Create(".million-timer")
	defer w.Close()
	b, _ := c.json.EncodePretty()
	w.Write(b)
}

type CacheRedis struct {
	conn redis.Conn
}

func NewCacheRedis(address string) *CacheRedis {
	c, err := redis.DialURL(address)
	if err != nil {
		panic(err)
	}

	return &CacheRedis{conn: c}
}

func (c *CacheRedis) GetBool(key string) bool {
	val, err := redis.Bool(c.conn.Do("GET", key))
	if err != nil {
		return false
	}
	return val
}

func (c *CacheRedis) GetString(key string) string {
	val, err := redis.String(c.conn.Do("GET", key))
	if err != nil {
		return ""
	}
	return val
}

func (c *CacheRedis) GetMap(key string) map[string]int {
	val, err := redis.IntMap(c.conn.Do("HGETALL", key))
	if err != nil {
		panic(err)
		return make(map[string]int)
	}
	return val
}

func (c *CacheRedis) Set(key string, val interface{}) {
	switch val.(type) {
	case map[string]int:
		c.conn.Do("HMSET", redis.Args{}.Add(key).AddFlat(val)...)
	default:
		c.conn.Do("SET", key, val)
	}
}

func (c *CacheRedis) Del(key string) {
	c.conn.Do("DEL", key)
}

func (c *CacheRedis) Close() {
	c.conn.Close()
}
