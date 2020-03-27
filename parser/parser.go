package parser

import (
	am "bitbucket.org/mikelsr/gauzaez/lexer/automaton"
	"github.com/mikelsr/bspl/proto"
)

const (
	// Tokens
	arrow        = "arrow"
	closeBrace   = "close_brace"
	closeBracket = "close_bracket"
	colon        = "colon"
	comma        = "comma"
	newline      = "newline"
	openBrace    = "open_brace"
	openBracket  = "open_bracket"
	whitespace   = "whitespace"
	word         = "word"

	// Reserved words

	// In input parameter
	In = "in"
	// Key parameter
	Key = "key"
	// Nil parameter of undefined scope
	Nil = "nil"
	// Out output parameter
	Out = "out"
	// Param parameter section declaration
	Param = "parameter"
	// Role declaration
	Role = "role"
)

var (
	// reservedWords contains every word reserved by BSPL
	reservedWords = []string{In, Key, Nil, Out, Param, Role}
)

// ProtoBuilder is used to parse a BSPL file and produce a protocol
type ProtoBuilder struct {
	p proto.Protocol
}

func (b *ProtoBuilder) parseName(tokens []am.Token, values []string) (int, bool) {
	i := nextNewline(tokens)
	// nextNewline should have encountered: [word, {, \n]
	if i == -1 || i != 2 {
		return 0, false
	}
	buff := struct {
		t []am.Token
		v []string
	}{t: tokens[:i], v: values[:i]}
	if string(buff.t[0]) != word && string(buff.t[1]) != openBracket {
		return 0, false
	}
	name := values[0]
	if isReserved(name) {
		return 0, false
	}
	b.p.Name = name
	return i, true
}

func (b *ProtoBuilder) parseRoles(tokens []am.Token, values []string) (int, bool) {
	i := nextNewline(tokens)
	// Expected at least role <Role>
	// Pair number: role <Role> <comma> <Role>...
	if i < 2 || i%2 != 0 {
		return 0, false
	}
	buff := struct {
		t []am.Token
		v []string
	}{t: tokens[:i], v: values[:i]}
	if buff.t[0] != word || buff.v[0] != Role {
		return 0, false
	}
	// Check validity of roles and separating commas
	roles := []proto.Role{}
	for j := 1; j < i; j++ {
		// even tokens are commas
		if j%2 == 0 {
			if buff.t[j] != comma {
				return 0, false
			}
		} else { // odd tokens are roles
			if buff.t[j] != word || isReserved(buff.v[j]) {
				return 0, false
			}
			r := proto.Role(buff.v[j])
			for _, v := range roles {
				if r == v {
					// repeated role
					return 0, false
				}
			}
			roles = append(roles, r)
		}
	}
	b.p.Roles = roles
	return i, true
}
