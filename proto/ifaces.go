package proto

import "reflect"

// Parameterizable  structures such as protocol or action definitions and implementations
type Parameterizable interface {
	Parameters() []Parameter
}

// Parameterized Parameterizable struct, having given values to parameters
type Parameterized interface {
	Value(Parameter) string
}

// Unique instance or something: a protocol instance or message
type Unique interface {
	Keys() map[string]string
}

// Differ between two Uniques: true if they are different, false if not
func Differ(u1, u2 Unique) bool {
	return !reflect.DeepEqual(u1.Keys(), u2.Keys())
}

// FindByName is a utility function to find a parameter by name in a slice
// of Parameters
func FindByName(name string, params []Parameter) (Parameter, bool) {
	for _, param := range params {
		if param.Name == name {
			return param, true
		}
	}
	return Parameter{}, false
}
