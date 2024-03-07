package pokeCache

import (
	"fmt"
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("testdata"),
		},
		{
			key: "https://example.com/path",
			val: []byte("moretestdata"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Add Get Test Case %v", i), func(t *testing.T) {
			cache := NewCache(interval)
			cache.Add(c.key, c.val)
			val, err := cache.Get(c.key)
			if err != nil {
				t.Errorf("expected to find key")
				return
			}
			if string(val) != string(c.val) {
				t.Errorf("expected to find value")
				return
			}
		})
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("Duplicate add Test Case %v", i), func(t *testing.T) {
			cache := NewCache(interval)
			cache.Add(c.key, c.val)
			err := cache.Add(c.key, []byte("extradata"))
			if err == nil {
				t.Errorf("expected error to be raised")
				return
			}
		})
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("Empty Get Test Case %v", i), func(t *testing.T) {
			cache := NewCache(interval)
			_, err := cache.Get(c.key)
			if err == nil {
				t.Errorf("expected error to be raised")
				return
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime * 2
	cache := NewCache(baseTime)
	cache.Add("https://example.com", []byte("testdata"))
	t.Run("Before reap", func(t *testing.T) {
		_, err := cache.Get("https://example.com")
		if err != nil {
			t.Errorf("expected to find key")
			return
		}
	})
	time.Sleep(waitTime)
	t.Run("After reap", func(t *testing.T) {
		_, err := cache.Get("https://example.com")
		if err == nil {
			t.Errorf("expected error to be raised")
			return
		}
	})

}
