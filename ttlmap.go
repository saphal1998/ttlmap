package ttlmap

import (
	"sync"
	"time"
)

type Hash int

type TTLMap interface {
	Get(key Hash) (interface{}, error)
	Add(key Hash, value interface{}, alive time.Duration) time.Duration
	Remove(key Hash) error
}

type mapValue struct {
	value interface{}
	done  <-chan time.Time
	hash  Hash
}

func (m *mapValue) suicide(t *ttlMap) {
	select {
	case <-m.done:
		{
			go t.Remove(m.hash)
		}
	}
}

type ttlMap struct {
	kvstore map[Hash]*mapValue
	mutex   *sync.RWMutex
}

func (t *ttlMap) Add(key Hash, value interface{}, alive time.Duration) time.Duration {
	t.mutex.RLock()
	_, ok := t.kvstore[key]
	t.mutex.RUnlock()
	if ok {
		t.Remove(key)
	}

	val := mapValue{
		value: value,
		hash:  key,
	}

	t.mutex.Lock()
	done := time.After(alive)
	val.done = done
	t.kvstore[key] = &val
	go val.suicide(t)
	t.mutex.Unlock()

	return alive
}

func (t *ttlMap) Remove(key Hash) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	_, ok := t.kvstore[key]
	if !ok {
		return ErrDoesNotExist
	}
	delete(t.kvstore, key)
	return nil
}

func (t *ttlMap) Get(key Hash) (interface{}, error) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	current, ok := t.kvstore[key]
	if !ok {
		return nil, ErrDoesNotExist
	}
	return current.value, nil
}

func NewTTLCache() TTLMap {
	return &ttlMap{
		kvstore: make(map[Hash]*mapValue),
		mutex:   &sync.RWMutex{},
	}
}
