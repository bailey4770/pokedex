package pokecache

import (
	"fmt"
	"testing"
	"time"
)

const interval = 2 * time.Second

var cases = []struct {
	key string
	val []byte
}{
	{
		key: "https://test.com",
		val: []byte("testdata"),
	},
	{
		key: "https://anothertest.co.uk/example",
		val: []byte("moretestdata"),
	},
	{
		key: "https://hey.co.uk/example",
		val: []byte("hello"),
	},
	{
		key: "https://no.co.uk/example",
		val: []byte("yes"),
	},
	{
		key: "https://big.co.uk/example",
		val: []byte("big"),
	},
	{
		key: "https://small.co.uk/example",
		val: []byte("small"),
	},
	{
		key: "https://german.co.uk/example",
		val: []byte("deutsch"),
	},
	{
		key: "https://french.co.uk/example",
		val: []byte("francais"),
	},
	{
		key: "https://lol.co.uk/example",
		val: []byte("wow"),
	},
}

func TestAddGet(t *testing.T) {
	for i, c := range cases {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			cache := NewCache(interval)

			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)

			if !ok {
				t.Errorf("expected to find key %q in cache, but it was not found", c.key)
				return
			}
			if string(val) != string(c.val) {
				t.Errorf("expected value %q for key %q, but got %q", c.val, c.key, val)
				return
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	t.Run("Filling and testing reapLoop", func(t *testing.T) {
		cache := NewCache(interval)

		for _, c := range cases {
			cache.Add(c.key, c.val)
		}

		time.Sleep(interval + 500*time.Millisecond)

		for _, c := range cases {
			val, ok := cache.Get(c.key)
			if ok {
				t.Errorf("expected key %q to be reaped from cache, but it was found with value %q", c.key, val)
				return
			}
		}
	})
}
