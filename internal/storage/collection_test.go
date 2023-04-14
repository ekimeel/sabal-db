package storage

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCollectionUInt32SerialKey(t *testing.T) {

	c := NewCollection[uint32, string]("test", 0, NewSerialKeyGenerator(0))

	for i := 1; i <= 100; i++ {
		key, err := c.Create(fmt.Sprintf("%d", i))
		assert.Equal(t, uint32(i), key)
		if err != nil {
			t.Fail()
		}
	}
	assert.Equal(t, 100, c.Len())
	for i := 1; i <= 100; i++ {
		v, loaded := c.Get(uint32(i))
		assert.True(t, loaded)
		assert.Equal(t, fmt.Sprintf("%d", i), v)
	}

}
