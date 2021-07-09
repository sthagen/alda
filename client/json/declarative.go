package json

import (
	"github.com/Jeffail/gabs/v2"
)

// Array is the JSON array constructor that gabs is missing. It's a function
// that takes 0 or more elements and returns a *gabs.Container representing a
// JSON array containing those elements.
//
// gabs is the only Go library I could find that comes close to correctly
// modeling JSON by using types and allowing you to construct and manage JSON
// objects and arrays programmatically.
//
// It's quirky, though, in that there doesn't appear to be any way to construct
// an array other than to use `New` to construct an object where one of the
// field values is an array, and then use `Search` to obtain a reference to the
// array.
//
// NOTE: There is gabs.Wrap, which takes any value and returns a *gabs.Container
// representing the JSON version of that value. That almost behaves correctly if
// you pass it a slice, but (probably because of Go's zero value semantics where
// an empty slice is equivalent to nil...) if you pass it an empty slice, it
// serializes the JSON value as null instead of [], which is just wrong.
func Array(elements ...interface{}) *Container {
	wrapperObject := gabs.New()
	wrapperObject.Array("theArray")

	for _, element := range elements {
		wrapperObject.ArrayAppend(element, "theArray")
	}

	return wrapperObject.Search("theArray")
}

// Object is gabs' missing declarative constructor for a JSON object.
func Object(elements ...interface{}) *Container {
	// I don't expect to ever use this function in an "at runtime" context,
	// otherwise I would make it so that the return type is (*Container, error).
	if len(elements)%2 != 0 {
		panic("The number of arguments must be even.")
	}

	object := gabs.New()

	for i := 0; i+1 < len(elements); i += 2 {
		k := elements[i]
		v := elements[i+1]
		// I don't expect to ever use this function in an "at runtime" context,
		// otherwise I would return an error if `k` isn't a string.
		object.Set(v, k.(string))
	}

	return object
}

// ToJson converts the provided value into a JSON container.
//
// This is an alias for the Wrap function in github.com/Jeffail/gabs/v2.
func ToJson(value interface{}) *Container {
	return gabs.Wrap(value)
}
