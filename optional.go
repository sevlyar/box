package box

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

// Optional represents optional value of type T.
// Optional value must be [Some] (i.e. having a value) or [None] (i.e. doesn't have a value).
// Optional is a comparable value-type. Don't recommend to use with big or complex types.
// Zero value of Optional is [None].
//
// Optional implements [sql.Scanner] and [driver.Valuer] interfaces.
// To database value conversion works, T should be one of the types accepted by [driver.Value]
// or implements the interfaces. [None] in database presented as NULL.
//
// Optional implements (un)marshalling from/to JSON. [None] value don't have presentation in JSON,
// so structure fields of type Optional should be annotated by `json:",omitzero"`
// to skip fields with the value. Calling of marshalling methods for [None] value will cause panic.
type Optional[T any] struct {
	some bool
	v    T
}

// Some returns [Optional] with the given value.
func Some[T any](val T) Optional[T] {
	return Optional[T]{
		some: true,
		v:    val,
	}
}

// None returns [None] value of [Optional]. It means no value.
func None[T any]() Optional[T] {
	return Optional[T]{}
}

// IsSome returns true if the [Optional] value is [Some].
func (opt Optional[T]) IsSome() bool {
	return opt.some
}

// IsNone returns true if the [Optional] value is [None].
func (opt Optional[T]) IsNone() bool {
	return !opt.some
}

// Get returns underlying value if [Optional] is [Some].
// Panics in case [Optional] is [None].
func (opt Optional[T]) Get() T {
	if !opt.some {
		panic("value is not presented")
	}

	return opt.v
}

// ToNullable returns an equivalent value of [Nullable] type.
func (opt Optional[T]) ToNullable() Nullable[T] {
	return Nullable[T]{
		valid: opt.some,
		v:     opt.v,
	}
}

var (
	_ driver.Valuer = Optional[any]{}
	_ sql.Scanner   = (*Optional[any])(nil)

	_ json.Marshaler   = Optional[any]{}
	_ json.Unmarshaler = (*Optional[any])(nil)
)

func (opt Optional[T]) Value() (driver.Value, error) {
	n := sql.Null[T]{
		Valid: opt.some,
		V:     opt.v,
	}

	return n.Value()
}

func (opt *Optional[T]) Scan(src any) error {
	var n sql.Null[T]

	if err := n.Scan(src); err != nil {
		return err
	}

	opt.some = n.Valid
	opt.v = n.V

	return nil
}

func (opt Optional[T]) IsZero() bool {
	return !opt.some
}

func (opt Optional[T]) MarshalJSON() ([]byte, error) {
	if !opt.some {
		panic("unable to marshal zero Optional[T] to JSON, use `json:\",omitzero\"` annotation for struct fields")
	}

	return marshalJSON(opt.v)
}

func (opt *Optional[T]) UnmarshalJSON(data []byte) error {
	err := unmarshalJSON(data, &opt.v)
	opt.some = err == nil

	return nil
}
