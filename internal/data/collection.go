package data

import (
	"errors"
	"sync"
)

const (
	errKeyViolation = "duplicate key violation"
	errEmptyKey     = "empty key"
	errKeyNotFound  = "key not found"
)

type Collection[V interface{}] struct {
	data map[string]V
	sync.RWMutex
}

func NewCollection[V interface{}](size int) Collection[V] {
	c := Collection[V]{}
	c.data = make(map[string]V, size)
	return c
}

func (c *Collection[V]) Get(key string) (V, bool) {
	c.RWMutex.RLock()
	defer c.RUnlock()
	val, ok := c.data[key]
	return val, ok
}

func (c *Collection[V]) Has(val string) bool {
	c.RWMutex.RLock()
	defer c.RUnlock()
	_, ok := c.data[val]
	return ok
}

func (c *Collection[V]) Len() int {
	return len(c.data)
}

func (c *Collection[V]) validateKey(key string) error {
	if len(key) == 0 {
		return errors.New(errEmptyKey)
	}
	return nil
}

func (c *Collection[V]) Create(key string, value V) error {

	err := c.validateKey(key)
	if err != nil {
		return err
	}

	c.RWMutex.Lock()
	defer c.Unlock()

	_, exists := c.data[key]
	if exists == false {
		c.data[key] = value
	} else {
		return errors.New(errKeyViolation)
	}

	return nil
}

func (c *Collection[V]) Replace(key string, value V) error {

	if c.Has(key) == false {
		return errors.New(errKeyNotFound)
	}

	c.RWMutex.Lock()
	c.data[key] = value
	defer c.Unlock()

	return nil
}

func (c *Collection[V]) Remove(key string) error {
	c.RWMutex.Lock()
	defer c.Unlock()
	delete(c.data, key)

	return nil
}

func (c *Collection[V]) Values() []V {
	values := make([]V, c.Len())

	c.RWMutex.RLock()
	defer c.RUnlock()
	var i int
	for _, v := range c.data {
		values[i] = v
		i++
	}
	return values
}

func (c *Collection[V]) Empty() bool {
	return len(c.data) == 0
}

func (c *Collection[V]) Filter(f func(string, V) bool) []V {

	filtered := make([]V, 0)
	c.RWMutex.RLock()
	for k, v := range c.data {
		if f(k, v) {
			filtered = append(filtered, v)
		}
	}
	defer c.RUnlock()

	return filtered
}
