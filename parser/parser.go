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
	// scopeWords contains keywords describing parameter scopes
	scopeWords = []string{In, Nil, Out}
)

// ProtoBuilder is used to parse a BSPL file and produce a protocol
type ProtoBuilder struct {
	p proto.Protocol
}

// parseName parses the protocol name declaration section of a BSPL protocol
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

// parseRoles parses the role declaration section of a BSPL protocol
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

// groupParamTokens parses a token slice and creates groups assuming
// the tokens belong to a parameter declaration
func groupParamTokens(tokens []am.Token, values []string) ([][]string, bool) {
	groups := make([][]string, 1)
	i := 0
	for j, t := range tokens {
		if t == comma {
			i++
			groups = append(groups, []string{})
			continue
		}
		if t == word {
			groups[i] = append(groups[i], values[j])
			continue
		}
		return [][]string{}, false
	}
	return groups, true
}

// parseParams extracts parameter from a token slice from a parameter declaration
// the expected token input is: <?scope> <name> <?key>, <?scope> <name> <?key>...
func parseParams(tokens []am.Token, values []string) ([]proto.Parameter, bool) {
	NIL := []proto.Parameter{}
	groups, ok := groupParamTokens(tokens, values)
	if !ok {
		return NIL, false
	}
	params := []proto.Parameter{}
	for _, g := range groups {
		if len(g) < 1 || len(g) > 3 {
			return NIL, false
		}
		scope := proto.Nil
		key := false
		var name string
		switch len(g) {
		// non-key, nil param
		case 1:
			name = g[0]
			if isReserved(name) {
				return NIL, false
			}
		// two cases: <param><key> and <scope><param>
		case 2:
			// case 1: <scope> <param>
			if g[1] == Key {
				name = g[0]
				key = true
				if isReserved(name) {
					return NIL, false
				}
			} else {
				if isReserved(g[1]) || !isScope(g[0]) {
					return NIL, false
				}
				scope = proto.IO(g[0])
				name = g[1]
			}
		case 3:
			if !isScope(g[0]) || isReserved(g[1]) || g[2] != Key {
				return NIL, false
			}
			scope = proto.IO(g[0])
			name = g[1]
			key = true
		default:
			return NIL, false
		}
		// check that the name is not repeated
		for _, p := range params {
			if p.Name == name {
				return NIL, false
			}
		}
		params = append(params, proto.Parameter{
			Io:   scope,
			Key:  key,
			Name: name,
		})
	}
	return params, true
}

func (b *ProtoBuilder) parseProtoParams(tokens []am.Token, values string) (int, bool) {
	return 0, false
}
