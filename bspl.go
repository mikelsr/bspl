package bspl

import (
	"io"

	"github.com/mikelsr/bspl/parser"
	"github.com/mikelsr/bspl/proto"
)

// Parse a BSPL protocol
func Parse(in io.Reader) (proto.Protocol, error) {
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
