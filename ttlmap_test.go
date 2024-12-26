package ttlmap

import (
	"testing"
	"time"
)

func TestAddValuesToCacheWithTimeout(t *testing.T) {
	key := Hash(1)
	expectedValue := "item1"
	expiry := time.Second

	cache := NewTTLCache()
	duration := cache.Add(key, expectedValue, expiry)
	if duration != expiry {
		t.Errorf("expected=%v, got=%v", expiry, duration)
	}

	gotValue, err := cache.Get(key)

	if err != nil || gotValue != expectedValue {
		t.Errorf("expectedValue=%v, but got error=%v, gotValue=%v", expectedValue, err, gotValue)
	}

	delta := 100 * time.Millisecond
	time.Sleep(expiry + delta)

	gotValue, err = cache.Get(key)

	if err == nil {
		t.Errorf("expectedErr=%v, but got err=%v, gotValue=%v", ErrDoesNotExist, err, gotValue)
	}
}

func TestAddValueAndRemoveToNotFind(t *testing.T) {
	key := Hash(1)
	expectedValue := "item1"
	expiry := time.Hour

	cache := NewTTLCache()
	cache.Add(key, expectedValue, expiry)
	gotValue, err := cache.Get(key)

	if err != nil || gotValue != expectedValue {
		t.Errorf("expectedValue=%v, but got error=%v, gotValue=%v", expectedValue, err, gotValue)
	}

	time.Sleep(time.Second)

	err = cache.Remove(key)
	if err != nil {
		t.Errorf("Expected value to be removed, but got=%v", err)

	}

	gotValue, err = cache.Get(key)

	if err == nil {
		t.Errorf("expectedErr=%v, but got err=nil, gotValue=%v", ErrDoesNotExist, gotValue)
	}
}

func TestAddingTheSameValue(t *testing.T) {
	key := Hash(1)
	expectedValue := "item1"
	highExpiry := 2 * time.Second
	lowExpiry := time.Second

	cache := NewTTLCache()
	cache.Add(key, expectedValue, highExpiry)
	gotValue, err := cache.Get(key)

	if err != nil || gotValue != expectedValue {
		t.Errorf("expectedValue=%v, but got error=%v, gotValue=%v", expectedValue, err, gotValue)
	}

	cache.Add(key, expectedValue, lowExpiry)
	gotValue, err = cache.Get(key)

	if err != nil || gotValue != expectedValue {
		t.Errorf("expectedValue=%v, but got error=%v, gotValue=%v", expectedValue, err, gotValue)
	}
	delta := 100 * time.Millisecond
	time.Sleep(lowExpiry + delta)

	gotValue, err = cache.Get(key)
	if err == nil {
		t.Errorf("expectedErr=%v, but got err=nil, gotValue=%v", ErrDoesNotExist, gotValue)
	}

	time.Sleep(highExpiry * 2)
}
