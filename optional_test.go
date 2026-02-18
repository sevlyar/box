package box

import (
	"encoding/json"
	"fmt"
)

// Zero value of Optional type is None.
// Use Some and None functions to construct value of Optional type.
func ExampleOptional_constructor() {
	var none1 Optional[string]
	none2 := None[string]()
	some := Some("value")

	fmt.Println(
		none1.IsNone(),
		none2.IsNone(),
		some.IsSome(),
	)
	// Output:
	// true true true
}

// Optional is a comparable value-type.
func ExampleOptional_compare() {
	fmt.Println(
		Some("a") == Some("a"),
		None[string]() == None[string](),
		Some("a") != Some("b"),
		Some("a") != None[string](),
	)
	// Output:
	// true true true true
}

// Get method allows you to get the underlying value,
// but checks whether the value is present and panics if it is not.
// So you should protect the Get method call with check the value is Some.
func ExampleOptional_check() {
	list := []Optional[int]{
		Some(1),
		None[int](),
		Some(2),
	}

	for _, opt := range list {
		if opt.IsSome() {
			v := opt.Get()
			fmt.Println(v)
		}
	}
	// Output:
	// 1
	// 2
}

// None presented in JSON as null. Use annotation `json:",omitzero"`
// to omit fields with None value from the encoding.
func ExampleOptional_marshalling() {
	type User struct {
		FirstName  string
		LastName   string
		MiddleName Optional[string] `json:",omitzero"`
		Age        Optional[int]
	}

	u := User{
		FirstName: "John",
		LastName:  "Doe",
	}

	b, _ := json.MarshalIndent(&u, "", "  ")

	fmt.Println(string(b))
	// Output:
	// {
	//   "FirstName": "John",
	//   "LastName": "Doe",
	//   "Age": null
	// }
}

// Use Optional2 type to work with twice optional values,
// when you need to build a modifying form of object with optional fields.
func ExampleOptional2_use() {
	var form struct {
		A Optional2[string]
		B Optional2[string]
		C Optional2[string]
	}

	input := `
	{
		"A": "str",
		"B": null
	}`

	_ = json.Unmarshal([]byte(input), &form)

	fmt.Println(
		form.A == Some2(Some("str")),
		form.B == Some2(None[string]()),
		form.C == None2[string](),
	)
	// Output:
	// true true true
}
