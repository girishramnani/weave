package orm

import (
	"reflect"

	"github.com/iov-one/weave/errors"
)

type limitedIterator struct {
	// n is count of remaining elements
	n int64
	// iterator is the underlying iterator
	iter SerialModelIterator
}

var _ SerialModelIterator = (*limitedIterator)(nil)

func (l *limitedIterator) LoadNext(dest SerialModel) error {
	if l.n > 0 {
		l.n--
		return l.LoadNext(dest)
	}
	return errors.Wrap(ErrLimit, "iterator limit reached")
}

func (l *limitedIterator) Release() {
	l.iter.Release()
}

func WithLimit(iter SerialModelIterator, limit int64) (SerialModelIterator, error) {
	if limit < 1 {
		return nil, errors.Wrap(errors.ErrInput, "invalid limit")
	}
	return &limitedIterator{iter: iter, n: limit}, nil
}

func ToSlice(iter SerialModelIterator, destination SerialModelSlicePtr) error {
	dest := reflect.ValueOf(destination)
	if dest.Kind() != reflect.Ptr {
		return errors.Wrap(errors.ErrType, "destination must be a pointer to slice of SerialModels")
	}
	if dest.IsNil() {
		return errors.Wrap(errors.ErrImmutable, "got nil pointer")
	}
	dest = dest.Elem()
	if dest.Kind() != reflect.Slice {
		return errors.Wrap(errors.ErrType, "destination must be a pointer to slice of SerialModels")
	}

	// Allocate memory depending on the type of the destination
	var elemTyp reflect.Type
	if dest.Type().Elem().Kind() != reflect.Ptr {
		elemTyp = dest.Type().Elem()
	} else {
		elemTyp = dest.Type().Elem().Elem()
	}
	d := reflect.New(elemTyp).Interface().(SerialModel)

	// Consume all the elements in iter to the destination
	for {
		if err := iter.LoadNext(d); err != nil {
			if !errors.ErrIteratorDone.Is(err) {
				return err
			}
			break
		}
		dest.Set(reflect.Append(dest, reflect.ValueOf(d)))
	}

	return nil
}
