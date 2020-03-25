package proto

// Action in a BSPL protocol
type Action struct {
	From       Role
	To         Role
	parameters []Parameter
}

// Parameters of an Action
func (a Action) Parameters() []Parameter {
	return a.parameters
}

// Parameter of a BSPL protocol
type Parameter struct {
	Io   IO
	Name string
	Key  bool
}

// Protocol is a definition of a BSPQL protocol
type Protocol struct {
	from       Role
	to         Role
	parameters []Parameter
}

// Parameters of a Protocol
func (p Protocol) Parameters() []Parameter {
	return p.parameters
}

// Role of a participant in a BSPL protocol
type Role string

// IO states wheter a parameter is an input or output parameter
type IO string

const (
	// In defines a local scope
	In IO = "in"
	// Out defines a global scope
	Out IO = "out"
	// Nil defines a parameter missing from a protocol instance
	Nil IO = "nil"
)
