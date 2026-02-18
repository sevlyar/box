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

	fmt.Println(none1.IsNone())
	fmt.Println(none2.IsNone())
	fmt.Println(some.IsSome())
	// Output:
	// true
	// true
	// true
}

// Optional is a comparable value-type.
func ExampleOptional_compare() {
	fmt.Println(Some("a") == Some("a"))
	fmt.Println(None[string]() == None[string]())
	fmt.Println(Some("a") != Some("b"))
	fmt.Println(Some("a") != None[string]())
	// Output:
	// true
	// true
	// true
	// true
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

// Use `json:",omitzero"` annotation for structure fields of type Optional to marshal None value properly.
func ExampleOptional_marshalling() {
	type User struct {
		FirstName  string
		LastName   string
		MiddleName Optional[string] `json:",omitzero"`
		Age        Optional[int]    `json:",omitzero"`
	}

	u := User{
		FirstName: "John",
		LastName:  "Doe",
		Age:       Some(18),
	}

	b, _ := json.MarshalIndent(&u, "", "  ")
	fmt.Println(string(b))
	// Output:
	// {
	//   "FirstName": "John",
	//   "LastName": "Doe",
	//   "Age": 18
	// }
}
