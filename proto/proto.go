package proto

// Action in a BSPL protocol
type Action struct {
	From       Role
	To         Role
	Parameters []Parameter
}

// IO states wheter a parameter is an input or output parameter
type IO string

// Parameter of a BSPL protocol
type Parameter struct {
	key  bool
	name string
	io   IO
}

// Protocol represents a BSPL protocol
type Protocol struct {
	Roles      []Role
	Parameters []Parameter
	Actions    []Action
}

// Role of a participant in a BSPL protocol
type Role string

const (
	// In defines a local scope
	In IO = "in"
	// Out defines a global scope
	Out IO = "out"
	// Nil defines a parameter missing from a protocol instance
	Nil IO = "nil"
)
