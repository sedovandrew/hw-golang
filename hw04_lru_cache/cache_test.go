package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(5)
		data := [...]Key{"zero", "one", "two", "three", "four"}
		for index, item := range data {
			c.Set(item, index)
		}

		c.Clear()

		for _, item := range data {
			value, ok := c.Get(item)
			require.False(t, ok)
			require.Nil(t, value)
		}
	})
}

func TestCacheOverflow(t *testing.T) {
	t.Run("simple cache overflow", func(t *testing.T) {
		c := NewCache(2)
		c.Set("first", 1)
		c.Set("second", 2)
		c.Set("third", 3)

		first, ok := c.Get("first")
		require.False(t, ok)
		require.Nil(t, first)
	})

	t.Run("complex cache overflow", func(t *testing.T) {
		testData := map[string]string{
			"one":   "first",
			"two":   "second",
			"three": "third",
			"four":  "fourth",
		}

		// Creating a cache with test data.
		c := NewCache(3)
		c.Set("one", testData["one"])     // [one]
		c.Set("two", testData["two"])     // [two, one]
		c.Set("three", testData["three"]) // [three, two, one]
		c.Get("one")                      // [one, three, two]
		c.Get("three")                    // [three, one, two]
		c.Set("four", testData["four"])   // [four, three, one]
		c.Get("two")                      // [four, three, one]
		c.Get("one")                      // [one, four, three]

		// Checking existing elements
		oneValue, ok := c.Get("one")
		require.True(t, ok)
		require.Equal(t, oneValue, testData["one"])

		threeValue, ok := c.Get("three")
		require.True(t, ok)
		require.Equal(t, threeValue, testData["three"])

		fourValue, ok := c.Get("four")
		require.True(t, ok)
		require.Equal(t, fourValue, testData["four"])

		// Checking missing elements
		twoValue, ok := c.Get("two")
		require.False(t, ok)
		require.Nil(t, twoValue)
	})
}

func TestCacheMultithreading(t *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
