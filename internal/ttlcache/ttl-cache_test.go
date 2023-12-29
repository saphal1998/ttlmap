package ttlcache

import (
	"testing"
	"time"
)

func TestSuccessCase(t *testing.T) {
	got := 2 + 2
	want := 4

	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestNestedSuccessCase(t *testing.T) {
	t.Run("add two numbers", func(t *testing.T) {
		got := 2 + 2
		want := 4

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
	t.Run("subtract two numbers", func(t *testing.T) {
		got := 2 - 2
		want := 0

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

// TODO :how does one test time? for instance I want to check that the cache is being cleared when ttl is in the order of minutes, how can it be simulated
func TestTTLCacheBasic(t *testing.T) {
	cache := NewCache()

	cache.Add("1", 1, uint64(time.Microsecond))

	time.Sleep(time.Microsecond / 2)

	got := cache.Get("1")
	want := 1

	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	time.Sleep(time.Microsecond)

	// Now the value shouldn't be present as the time has expired
	got = cache.Get("1")

	if got != nil {
		t.Errorf("got %v, want nil", got)
	}
}
