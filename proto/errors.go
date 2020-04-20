package proto

import "fmt"

// ValidationError is created while validating protocols
type ValidationError struct {
	Err error
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("Validation error: \"%s\"", e.Err.Error())
}
