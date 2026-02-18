/*
Package box provides generic types to reduce memory footprint and GC pressure.

# Optional and Nullable

[Optional] and [Nullable] types are equivalent in most cases but there are several differences:
  - [Optional] presents optional value, and [None] means "no value", "no field", "no item";
  - [Nullable] just adds [Null] value to the underlying type values set.

So that there are some differences of (un)marshalling from/to JSON:

  - [Some]  - presented as underlying value
  - [None]  - causes panic on direct marshalling, but with annotation `json:",omitzero"` is able to hide field
  - [Valid] - presented as underlying type
  - [Null]  - presented as null

There is no difference in database type conversions between [Optional] and [Nullable] types.

Is it possible to combine types for complex scenarios. For example, JSON forms:

	type EditDocumentForm struct {
		// None means "field unmodified"
		// Null means "now field is unset"
		Title Optional[Nullable[string]]
		// ...
	}
*/
package box
