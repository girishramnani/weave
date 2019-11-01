package orm

import (
	"testing"

	"github.com/iov-one/weave/store"
	"github.com/iov-one/weave/weavetest/assert"
)

func TestToSlice(t *testing.T) {
	db := store.MemStore()

	b := NewSerialModelBucket("cnts", &CounterWithID{})

	cntr := CounterWithID{Count: 1}
	// make sure we point to value in array, so this PrimaryKey gets set
	err := b.Create(db, &cntr)
	assert.Nil(t, err)

	iter, err := b.IndexScan(db, "counter", nil, false)

	var dest1 []*CounterWithID
	err = ToSlice(iter, &dest1)
	assert.Nil(t, err)
	assert.Equal(t, cntr, dest1[0])

	var dest2 []CounterWithID
	err = ToSlice(iter, dest2)
	assert.Nil(t, err)
}
