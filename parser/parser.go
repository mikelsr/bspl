package parser

import (
	"errors"
	"fmt"

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
func (b *ProtoBuilder) parseName(tokens []am.Token, values []string) (int, error) {
	i := nextNewline(tokens)
	// nextNewline should have encountered: [word, {, \n]
	if i == -1 || i != 2 {
		return 0, fmt.Errorf("Expected (<protocol> {) got (%s)", values[:i])
	}
	buff := struct {
		t []am.Token
		v []string
	}{t: tokens[:i], v: values[:i]}
	if string(buff.t[0]) != word && string(buff.t[1]) != openBracket {
		return 0, ParseError{Expected: "protocol name or bracket", Found: fmt.Sprint(buff.v[:2])}
	}
	name := values[0]
	if isReserved(name) {
		return 0, ReservedError{Word: name}
	}
	b.p.Name = name
	return i, nil
}

// parseRoles parses the role declaration section of a BSPL protocol
func (b *ProtoBuilder) parseRoles(tokens []am.Token, values []string) (int, error) {
	i := nextNewline(tokens)
	// Expected at least role <Role>
	// Pair number: role <Role> <comma> <Role>...
	if i < 2 || i%2 != 0 {
		return 0, errors.New("Invalid parameter definition")
	}
	buff := struct {
		t []am.Token
		v []string
	}{t: tokens[:i], v: values[:i]}
	if buff.t[0] != word || buff.v[0] != Role {
		return 0, ParseError{Expected: "role", Found: buff.v[0]}
	}
	// Check validity of roles and separating commas
	roles := []proto.Role{}
	for j := 1; j < i; j++ {
		// even tokens are commas
		if j%2 == 0 {
			if buff.t[j] != comma {
				return 0, ParseError{Expected: ",", Found: buff.v[j]}
			}
		} else { // odd tokens are roles
			if buff.t[j] != word || isReserved(buff.v[j]) {
				return 0, ReservedError{Word: buff.v[j]}
			}
			r := proto.Role(buff.v[j])
			for _, v := range roles {
				if r == v {
					// repeated role
					return 0, fmt.Errorf("Repeated role: %s", r)
				}
			}
			roles = append(roles, r)
		}
	}
	b.p.Roles = roles
	return i, nil
}

// groupParamTokens parses a token slice and creates groups assuming
// the tokens belong to a parameter declarationz
func groupParamTokens(tokens []am.Token, values []string) ([][]string, error) {
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
		return [][]string{}, ParseError{Expected: "word or ','", Found: values[j]}
	}
	return groups, nil
}

// parseParams extracts parameter from a token slice from a parameter declaration
// the expected token input is: <?scope> <name> <?key>, <?scope> <name> <?key>...
func parseParams(tokens []am.Token, values []string) ([]proto.Parameter, error) {
	NIL := []proto.Parameter{}
	groups, err := groupParamTokens(tokens, values)
	if err != nil {
		return NIL, err
	}
	params := []proto.Parameter{}
	for _, g := range groups {
		if len(g) < 1 || len(g) > 3 {
			return NIL, ParamError{Comp: values}
		}
		scope := proto.Nil
		key := false
		var name string
		switch len(g) {
		// non-key, nil param
		case 1:
			name = g[0]
			if isReserved(name) {
				return NIL, ReservedError{Word: name}
			}
		// two cases: <param><key> and <scope><param>
		case 2:
			// case 1: <scope> <param>
			if g[1] == Key {
				name = g[0]
				key = true
				if isReserved(name) {
					return NIL, ReservedError{Word: name}
				}
			} else {
				if isReserved(g[1]) {
					return NIL, ReservedError{Word: name}
				} else if !isScope(g[0]) {
					return NIL, ParseError{Expected: scopeWords, Found: g[0]}
				}
				scope = proto.IO(g[0])
				name = g[1]
			}
		case 3:
			if !isScope(g[0]) || isReserved(g[1]) || g[2] != Key {
				return NIL, ParseError{Expected: "<?scope> <name> <?key>", Found: g}
			}
			scope = proto.IO(g[0])
			name = g[1]
			key = true
		}
		// check that the name is not repeated
		for _, p := range params {
			if p.Name == name {
				return NIL, fmt.Errorf("Repeated parameter name: %s", name)
			}
		}
		params = append(params, proto.Parameter{
			Io:   scope,
			Key:  key,
			Name: name,
		})
	}
	return params, nil
}

func (b *ProtoBuilder) parseProtoParams(tokens []am.Token, values string) (int, bool) {
	return 0, false
}
