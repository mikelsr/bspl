package bspl

import (
	"io"

	"github.com/mikelsr/bspl/parser"
	"github.com/mikelsr/bspl/proto"
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
	tokens, err := parser.LexStream(in)
	if err != nil {
		return proto.Protocol{}, err
	}
	stripped := parser.Strip(*tokens)
	b := new(parser.ProtoBuilder)
	if err := b.Parse(stripped.Tokens, stripped.Values); err != nil {
		return proto.Protocol{}, err
	}
	return b.Protocol(), nil
}
