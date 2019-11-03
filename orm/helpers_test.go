package orm

import (
	"reflect"
	"testing"

	"github.com/iov-one/weave/store"
	"github.com/iov-one/weave/weavetest/assert"
)

func TestWithLimit(t *testing.T) {
	db := store.MemStore()

	b := NewSerialModelBucket("cnts", &CounterWithID{},
		WithIndexSerial("counter", func(Object) ([]byte, error) { return []byte("all"), nil }, false))

	cntr1 := CounterWithID{Count: 1}
	// make sure we point to value in array, so this PrimaryKey gets set
	err := b.Create(db, &cntr1)
	assert.Nil(t, err)

	iter, err := b.IndexScan(db, "counter", nil, false)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	var dest []*CounterWithID
	err = ToSlice(iter, &dest)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	assert.Equal(t, cntr1, *dest[0])

	var dest2 []CounterWithID
	err = ToSlice(iter, &dest2)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}

func TestToSlice(t *testing.T) {
	db := store.MemStore()

	b := NewSerialModelBucket("cnts", &CounterWithID{},
		WithIndexSerial("counter", func(Object) ([]byte, error) { return []byte("all"), nil }, false))

	expected := []*CounterWithID{
		&CounterWithID{Count: 1},
		&CounterWithID{Count: 2},
	}

	for _, e := range expected {
		// make sure we point to value in array, so this PrimaryKey gets set
		err := b.Create(db, e)
		assert.Nil(t, err)
	}

	iter, err := b.IndexScan(db, "counter", nil, false)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	var dest []*CounterWithID
	err = ToSlice(iter, &dest)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !reflect.DeepEqual(dest, expected) {
		t.Errorf("values do not match, expected: %+v, got: %+v", expected, dest)
	}
}
