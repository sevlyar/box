package box

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding"
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
// Optional implements (un)marshalling from/to JSON. [None] value presented as null.
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

var nullStrBytes = []byte("null")

func (opt Optional[T]) MarshalJSON() ([]byte, error) {
	if !opt.some {
		return nullStrBytes, nil
	}

	return marshalJSON(opt.v)
}

func marshalJSON[T any](v T) ([]byte, error) {
	if casted, ok := any(v).(json.Marshaler); ok {
		return casted.MarshalJSON()
	}

	if casted, ok := any(v).(encoding.TextMarshaler); ok {
		return casted.MarshalText()
	}

	return json.Marshal(v)
}

func (opt *Optional[T]) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, nullStrBytes) {
		*opt = None[T]()
		return nil
	}

	err := unmarshalJSON(data, &opt.v)
	opt.some = err == nil

	return nil
}

func unmarshalJSON[T any](data []byte, ptr *T) error {
	if casted, ok := any(ptr).(json.Unmarshaler); ok {
		return casted.UnmarshalJSON(data)
	}

	if casted, ok := any(ptr).(encoding.TextUnmarshaler); ok {
		return casted.UnmarshalText(data)
	}

	return json.Unmarshal(data, ptr)
}

// Optional2 presents twice optional value: Optional[Optional[T]].
//
// Marshalling of Optional2 type to JSON is decorated. [None2] value can't marshalled to JSON directly.
// Calling of marshalling method for [None2] value will cause panic. So fields of type Optional2 should be
// annotated by `json:",omitzero"` to omit from the encoding when value is [None2].
type Optional2[T any] struct {
	Optional[Optional[T]]
}

func Some2[T any](val Optional[T]) Optional2[T] {
	return Optional2[T]{Some(val)}
}

func None2[T any]() Optional2[T] {
	return Optional2[T]{}
}

func (opt2 Optional2[T]) IsZero() bool {
	return opt2.IsNone()
}

func (opt2 Optional2[T]) MarshalJSON() ([]byte, error) {
	if opt2.IsNone() {
		panic("unable to marshal zero Optional2[T] to JSON, use `json:\",omitzero\"` annotation for struct fields")
	}

	return opt2.Optional.Get().MarshalJSON()
}

func (opt2 *Optional2[T]) UnmarshalJSON(data []byte) error {
	err := opt2.v.UnmarshalJSON(data)
	opt2.some = err == nil

	return nil
}
