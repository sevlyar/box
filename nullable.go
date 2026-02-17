package box

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

type Nullable[T any] struct {
	valid bool
	v     T
}

func Valid[T any](val T) Nullable[T] {
	return Nullable[T]{
		valid: true,
		v:     val,
	}
}

func Null[T any]() Nullable[T] {
	return Nullable[T]{}
}

func (n Nullable[T]) IsValid() bool {
	return n.valid
}

func (n Nullable[T]) IsNull() bool {
	return !n.valid
}

func (n Nullable[T]) Get() T {
	if !n.valid {
		panic("value is not presented")
	}

	return n.v
}

func (n Nullable[T]) ToOptional() Optional[T] {
	return Optional[T]{
		some: n.valid,
		v:    n.v,
	}
}

var (
	_ driver.Valuer = Nullable[any]{}
	_ sql.Scanner   = (*Nullable[any])(nil)

	_ json.Marshaler   = Nullable[any]{}
	_ json.Unmarshaler = (*Nullable[any])(nil)
)

func (n Nullable[T]) Value() (driver.Value, error) {
	sqlNull := sql.Null[T]{
		Valid: n.valid,
		V:     n.v,
	}

	return sqlNull.Value()
}

func (n *Nullable[T]) Scan(src any) error {
	var sqlNull sql.Null[T]

	if err := n.Scan(src); err != nil {
		return err
	}

	n.valid = sqlNull.Valid
	n.v = sqlNull.V

	return nil
}

var nullStrBytes = []byte("null")

func (n Nullable[T]) MarshalJSON() ([]byte, error) {
	if !n.valid {
		return nullStrBytes, nil
	}

	return marshalJSON(n.v)
}

func (n *Nullable[T]) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, nullStrBytes) {
		*n = Null[T]()
		return nil
	}

	err := unmarshalJSON(data, &n.v)
	n.valid = err == nil

	return err
}
