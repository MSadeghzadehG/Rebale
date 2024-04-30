package main

import (
	"testing"

	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRebaleConnect(t *testing.T) {
	t.Run("Test Connect func", func(t *testing.T) {
		var c Rebale
		err := c.Connect("127.0.0.1:6937")
		require.NoError(t, err)
	})
}

func TestRebalePing(t *testing.T) {
	t.Run("Test Ping func", func(t *testing.T) {
		var c Rebale
		err := c.Ping()
		require.NoError(t, err)
	})
}

func TestRebaleSet(t *testing.T) {
	t.Run("Test Set func", func(t *testing.T) {
		var c Rebale
		err := c.Connect("127.0.0.1:6937")
		require.NoError(t, err)
		cRedis := redis.NewClient(&redis.Options{PoolSize: 1})

		cases := map[string]string{
			"hello": "سلام",
			"empty": "",
		}
		for k, v := range cases {
			err = c.Set(k, v, len(v))
			require.NoError(t, err)
			rValue, err := cRedis.Get(k).Result()
			require.NoError(t, err)
			assert.EqualValues(t, rValue, v)
		}
	})
}

func TestRebaleGet(t *testing.T) {
	t.Run("Test Set func", func(t *testing.T) {
		var c Rebale
		err := c.Connect("127.0.0.1:6937")
		require.NoError(t, err)
		cRedis := redis.NewClient(&redis.Options{PoolSize: 1})

		cases := map[string]string{
			"hello": "سلام",
			"empty": "",
		}
		for k, v := range cases {
			err = c.Set(k, v, len(v))
			require.NoError(t, err)
			cValue, err := c.Get(k)
			assert.EqualValues(t, cValue, v)
			rValue, err := cRedis.Get(k).Result()
			require.NoError(t, err)
			assert.EqualValues(t, rValue, v)
		}
	})
}
