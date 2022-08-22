package directusapi

import (
	"encoding/json"
	"reflect"
)

type Optional[T any] struct {
	value T
	op    operation
}

func UnsetOptional[T any]() Optional[T] {
	return Optional[T]{
		op: unset,
	}
}

func SetOptional[T any](val T) Optional[T] {
	return Optional[T]{
		value: val,
		op:    set,
	}
}

func (o Optional[T]) ValueOrZero() T {
	if o.op != set {
		var zeroval T
		return zeroval
	}
	return o.value
}

func (o Optional[T]) ValueMust() T {
	if o.op != set {
		panic("value is not set")
	}
	return o.value
}

func (o Optional[T]) IsSet() bool {
	return o.op == set
}

func (o Optional[T]) MarshalJSON() ([]byte, error) {
	switch o.op {
	case set:
		return json.Marshal(o.value)
	case unset:
		return []byte(`null`), nil
	default:
		// https://github.com/golang/go/issues/11939
		return json.Marshal(nil)
	}
}

func (o *Optional[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		o.op = unset
		return nil
	}

	if err := json.Unmarshal(data, &o.value); err != nil {
		return err
	}
	o.op = set

	// in case of noop the UnmarhsalJSON is not called
	return nil
}

func (o Optional[T]) getOp() operation {
	// this is hack for reflection pkg
	return o.op
}

func (o Optional[T]) fields(prefix string) []string {
	var optVal T
	f := reflect.TypeOf(optVal)

	if f.Kind() == reflect.Struct {
		var t Time
		isTime := f.ConvertibleTo(reflect.TypeOf(t))
		isOptional := f.Implements(reflect.TypeOf(new(isOpt)).Elem())
		if isOptional {
			panic("optional of optional is not supported")
		}
		if isTime {
			return []string{prefix}
		}
		return iterateFields(f, prefix)
	}
	return []string{prefix}
}

// 1. don't touch the value
// 		=> zero value of Optional[T]
// 2. unset the value (null)
//  	=> UnsetOptional[int]()
// 3. set the value
//		=> SetOptional(3)

type operation uint8

const (
	noop operation = iota
	unset
	set
)
