package parser

import "fmt"

// ParamError is returned when a parameter is incorreclty declared
type ParamError struct {
	Comp []string
}

func (e ParamError) Error() string {
	return fmt.Sprintf("Invalid parameter definition: %s.", e.Comp)
}

// ParseError is returned when an unexpected character is found
type ParseError struct {
	Expected interface{}
	Found    interface{}
}

func (e ParseError) Error() string {
	return fmt.Sprintf("Expected characters '%s' but found '%s'.",
		e.Expected, e.Found)
}

// ReservedError is returned when a reserved BSPL word is used as a value
type ReservedError struct {
	Word string
}

func (e ReservedError) Error() string {
	return fmt.Sprintf("Used reserved word as a value, words reserved by BSPL are %s.",
		reservedWords)
}
