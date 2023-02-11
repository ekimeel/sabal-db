package data

import (
	"fmt"
	"github.com/ekimeel/db-api/pb"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCollection(t *testing.T) {
	c := NewCollection[pb.Point](0)
	assert.True(t, c.Empty())
}

func TestCollection_Get(t *testing.T) {
	c := NewCollection[*pb.Point](0)

	pointA := pb.Point{Uuid: "abc"}
	pointB := pb.Point{Uuid: "bcd"}
	pointC := pb.Point{Uuid: "cde"}

	c.Create(pointA.Uuid, &pointA)
	c.Create(pointB.Uuid, &pointB)
	c.Create(pointC.Uuid, &pointC)

	assert.Equal(t, c.Get(pointA.Uuid).Uuid, pointA.Uuid)
	assert.Equal(t, c.Get(pointB.Uuid).Uuid, pointB.Uuid)
	assert.Equal(t, c.Get(pointC.Uuid).Uuid, pointC.Uuid)
}

func TestCollection_Values(t *testing.T) {

	c := NewCollection[*pb.Point](0)

	for i := 0; i < 10000; i++ {
		uuid := fmt.Sprintf("%d", i)
		c.Create(uuid, &pb.Point{Uuid: uuid})
	}

	values := c.Values()
	assert.Equal(t, 10000, len(values))

}
