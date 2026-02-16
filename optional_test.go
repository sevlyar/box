package box

import "fmt"

// Zero value of Optional type is None.
// Use Some and None functions to construct value of Optional type.
func ExampleOptional_ctor() {
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
func ExampleOptional_handle() {
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
