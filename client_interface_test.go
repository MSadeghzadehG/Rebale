package main

import (
	"strings"
	"testing"

	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRebaleConnect(t *testing.T) {
	t.Run("Test Connect func", func(t *testing.T) {
		c := &MyRebaleImpl{}
		err := c.Connect("127.0.0.1:6379")
		require.NoError(t, err)
	})
}

func TestRebalePing(t *testing.T) {
	t.Run("Test Ping func", func(t *testing.T) {
		c := &MyRebaleImpl{}
		err := c.Ping()
		require.NoError(t, err)
	})
}

func TestRebaleSet(t *testing.T) {
	t.Run("Test Set func", func(t *testing.T) {
		c := &MyRebaleImpl{}
		err := c.Connect("127.0.0.1:6379")
		require.NoError(t, err)
		cRedis := redis.NewClient(&redis.Options{PoolSize: 1})

		cases := map[string]string{
			"hello": "سلام",
			"empty": "",
		}
		for k, v := range cases {
			err = c.Set(k, strings.NewReader(v), len(v))
			require.NoError(t, err)
			rValue, err := cRedis.Get(k).Result()
			require.NoError(t, err)
			assert.EqualValues(t, rValue, v)
		}
	})
}

func TestRebaleGet(t *testing.T) {
	t.Run("Test Set func", func(t *testing.T) {
		c := &MyRebaleImpl{}
		err := c.Connect("127.0.0.1:6379")
		require.NoError(t, err)
		cRedis := redis.NewClient(&redis.Options{PoolSize: 1})

		cases := map[string]string{
			"hello": "سلام",
			"empty": "",
		}
		for k, v := range cases {
			err = c.Set(k, strings.NewReader(v), len(v))
			require.NoError(t, err)
			cValue, err := c.Get(k)
			assert.Nil(t, err)
			assert.EqualValues(t, cValue, v)
			rValue, err := cRedis.Get(k).Result()
			require.NoError(t, err)
			assert.EqualValues(t, rValue, v)
		}
	})
}
