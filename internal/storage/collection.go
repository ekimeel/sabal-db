package storage

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"sync"
)

const (
	errKeyViolation = "duplicate key violation"
	errEmptyKey     = "empty key"
	errKeyNotFound  = "key not found"
)

type Collection[K comparable, V interface{}] struct {
	Name         string
	data         map[K]V
	keyGenerator KeyGenerator[K]
	sync.RWMutex
}

func NewCollection[K comparable, V interface{}](name string, size int, keyGenerator KeyGenerator[K]) *Collection[K, V] {
	log.Infof("creating new collection [%s]", name)
	c := &Collection[K, V]{}
	c.Name = name
	c.data = make(map[K]V, size)
	c.keyGenerator = keyGenerator
	return c
}

func (c *Collection[K, V]) Get(key K) (V, bool) {
	c.RWMutex.RLock()
	defer c.RUnlock()
	val, ok := c.data[key]
	return val, ok
}

func (c *Collection[K, V]) Has(key K) bool {
	c.RWMutex.RLock()
	defer c.RUnlock()
	_, ok := c.data[key]
	return ok
}

func (c *Collection[K, V]) Len() int {
	return len(c.data)
}

func (c *Collection[K, V]) Create(value V) (K, error) {
	c.RWMutex.Lock()

	key := c.keyGenerator.NextKey()
	_, exists := c.data[key]
	if exists == false {
		c.data[key] = value
	} else {
		return key, errors.New(errKeyViolation)
	}

	defer c.Unlock()
	return key, nil

}

func (c *Collection[K, V]) CreateWithKey(key K, value V) error {
	c.RWMutex.Lock()
	c.data[key] = value
	defer c.Unlock()
	return nil
}

func (c *Collection[K, V]) Replace(key K, value V) error {

	if c.Has(key) == false {
		return errors.New(errKeyNotFound)
	}

	c.RWMutex.Lock()
	c.data[key] = value
	defer c.Unlock()

	return nil
}

func (c *Collection[K, V]) Remove(key K) error {
	c.RWMutex.Lock()
	defer c.Unlock()
	delete(c.data, key)

	return nil
}

func (c *Collection[K, V]) Values() []V {
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

func (c *Collection[K, V]) Empty() bool {
	return len(c.data) == 0
}

func (c *Collection[K, V]) Filter(f func(K, V) bool) []V {

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
