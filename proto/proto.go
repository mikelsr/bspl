package proto

import "strings"

// Action in a BSPL protocol
type Action struct {
	Name   string
	From   Role
	To     Role
	Params []Parameter
}

// Parameters of an Action
func (a Action) Parameters() []Parameter {
	return a.Params
}

// Parameter of a BSPL protocol
type Parameter struct {
	Io   IO
	Key  bool
	Name string
}

// Protocol is a definition of a BSPQL protocol
type Protocol struct {
	Actions []Action
	Name    string
	Roles   []Role
	Params  []Parameter
}

// Parameters of a Protocol
func (p Protocol) Parameters() []Parameter {
	return p.Params
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

// Key of the protocol
func (p Protocol) Key() string {
	keys := p.Keys()
	n := len(keys)
	var sb strings.Builder

	sb.WriteString(p.Name)
	sb.WriteRune(KeySeparator)
	for i, k := range keys {
		sb.WriteString(k.Name)
		if i != n-1 {
			sb.WriteRune(KeySeparator)
		}
	}
	return sb.String()
}

// Sort all elements of a protocol
func (p *Protocol) Sort() {
	SortRoles(p.Roles)
	SortParameters(p.Params)
	for _, a := range p.Actions {
		SortParameters(a.Params)
	}
	SortActions(p.Actions)
}

func (a Action) String() string {
	var s strings.Builder
	s.WriteString(string(a.From) + " -> " + string(a.To) + ": " + a.Name + "[")
	if len(a.Params) > 0 {
		s.WriteString(a.Params[0].String())
		for _, p := range a.Params[1:] {
			s.WriteString(", " + p.String())
		}
	}
	s.WriteString("]")
	return s.String()
}

func (p Parameter) String() string {
	var s strings.Builder
	if p.Io != Nil {
		s.WriteString(string(p.Io) + " ")
	}
	s.WriteString(p.Name)
	if p.Key {
		s.WriteString(" key")
	}
	return s.String()
}

func (p Protocol) String() string {
	var s strings.Builder
	s.WriteString(p.Name + " {\n\trole ")
	s.WriteString(string(p.Roles[0]))
	for _, r := range p.Roles[1:] {
		s.WriteString(", " + string(r))
	}
	s.WriteString("\n\tparameter " + p.Params[0].String())
	for _, v := range p.Params[1:] {
		s.WriteString(", " + v.String())
	}
	s.WriteString("\n\n")
	for _, a := range p.Actions {
		s.WriteString("\t" + a.String() + "\n")
	}
	s.WriteString("}")
	return s.String()
}

func findKeys(params []Parameter) []Parameter {
	keyParams := make([]Parameter, 0)
	for _, param := range params {
		if param.Key {
			keyParams = append(keyParams, param)
		}
	}
	return keyParams
}

func findIns(params []Parameter) []Parameter {
	inParams := make([]Parameter, 0)
	for _, param := range params {
		if param.Io == In {
			inParams = append(inParams, param)
		}
	}
	return inParams
}

func findNils(params []Parameter) []Parameter {
	nilParams := make([]Parameter, 0)
	for _, param := range params {
		if param.Io == Nil {
			nilParams = append(nilParams, param)
		}
	}
	return nilParams
}

func findOuts(params []Parameter) []Parameter {
	outParams := make([]Parameter, 0)
	for _, param := range params {
		if param.Io == Out {
			outParams = append(outParams, param)
		}
	}
	return outParams
}

// Keys returns a list of the key parameters of the protocol
func (p Protocol) Keys() []Parameter {
	return findKeys(p.Parameters())
}

// Ins returns a list of the implicit parameters of the protocol
func (p Protocol) Ins() []Parameter {
	return findIns(p.Params)
}

// Outs returns a list of the explicit parameters of the protocol
func (p Protocol) Outs() []Parameter {
	return findOuts(p.Params)
}

// Nils returns a list of the nil parameters of the protocol
func (p Protocol) Nils() []Parameter {
	return findNils(p.Params)
}

// Keys returns a list of the key parameters of the action
func (a Action) Keys() []Parameter {
	return findKeys(a.Parameters())
}

// Ins returns a list of the implicit parameters of the action
func (a Action) Ins() []Parameter {
	return findIns(a.Params)
}

// Nils returns a list of the nil parameters of the action
func (a Action) Nils() []Parameter {
	return findNils(a.Params)
}

// Outs returns a list of the explicit parameters of the action
func (a Action) Outs() []Parameter {
	return findOuts(a.Params)
}
