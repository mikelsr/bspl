package bspl

import (
	"io"
	"reflect"

	"github.com/mikelsr/bspl/parser"
	"github.com/mikelsr/bspl/proto"
	"github.com/mikelsr/bspl/reason"
)

type (
	// Action is an alias for proto.Action
	Action = proto.Action
	// IO is an alias for proto.IO
	IO = proto.IO
	// Parameter is an alias for proto.Parameter
	Parameter = proto.Parameter
	// Protocol is an alias for proto.Protocol
	Protocol = proto.Protocol
	// Role is an alias for proto.Role
	Role = proto.Role

	// Reasoner is an alias for reason.Reasoner
	Reasoner = reason.Reasoner
	// Instance is an alias for reason.Instance
	Instance = reason.Instance
	// Message is an alias for reason.Message
	Message = reason.Message
	// Messages is an alias for reason.Messages
	Messages = reason.Messages
	// Roles is an alias for reason.Roles
	Roles = reason.Roles
	// Values is an alias for reason.Values
	Values = reason.Values
)

const (
	// In defines a local scope
	In IO = proto.In
	// Out defines a global scope
	Out IO = proto.Out
	// Nil defines a parameter missing from a protocol instance
	Nil IO = proto.Nil
)

// Parse a BSPL protocol
func Parse(in io.Reader) (Protocol, error) {
	return parser.Parse(in)
}

// Compare two BSPL protocols
func Compare(a, b Protocol) bool {
	return reflect.DeepEqual(a, b)
}
