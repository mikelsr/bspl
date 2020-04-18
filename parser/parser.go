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

// Protocol of the ProtoBuilder
func (b ProtoBuilder) Protocol() proto.Protocol {
	return b.p
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
		return 0, errors.New("Invalid role definition")
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

func (b *ProtoBuilder) parseProtoParams(tokens []am.Token, values []string) (int, error) {
	i := nextNewline(tokens)
	// minimal number of word tokens: [parameter, out, X, key] = 4
	if i < 4 {
		return 0, errors.New("Invalid protocol parameter definition")
	}
	buff := struct {
		t []am.Token
		v []string
	}{t: tokens[:i], v: values[:i]}
	// first word is "parameter"
	if buff.t[0] != word || buff.v[0] != Param {
		return 0, ParseError{Expected: Param, Found: buff.v[0]}
	}
	params, err := parseParams(buff.t[1:], buff.v[1:])
	if err != nil {
		return 0, err
	}
	// ensure that at least one parameter is a key parameter
	keyParam := false
	for _, p := range params {
		if p.Key {
			keyParam = true
			break
		}
	}
	if !keyParam {
		return 0, errors.New("No key parameters")
	}
	b.p.Params = params
	return i, nil
}

func (b *ProtoBuilder) parseActions(tokens []am.Token, values []string) (int, error) {
	i := nextNewline(tokens)
	var actionName string
	var from, to proto.Role
	// minimal number of tokens: RoleA -> RoleB: Act[P] = 8
	if i < 8 {
		return 0, errors.New("Invalid action")
	}
	buff := struct {
		t []am.Token
		v []string
	}{t: tokens[:i], v: values[:i]}

	// first and third tokens are Roles
	for _, i := range []int{0, 2} {
		if buff.t[i] != word {
			return 0, ParseError{Expected: "<Role>", Found: buff.v[i]}
		}
		// the role must match one of the roles declared previously
		definedRole := false
		for _, role := range b.p.Roles {
			if string(role) == buff.v[i] {
				definedRole = true
			}
		}
		if !definedRole {
			return 0, fmt.Errorf("Unknown role: %s", buff.v[i])
		}
	}
	from = proto.Role(buff.v[0])
	to = proto.Role(buff.v[2])

	if buff.t[1] != arrow {
		return 0, ParseError{Expected: "->", Found: buff.v[1]}
	}

	if buff.t[3] != colon {
		return 0, ParseError{Expected: ":", Found: buff.v[3]}
	}

	if buff.t[4] != word {
		return 0, ParseError{Expected: "<Action name>", Found: buff.v[4]}
	}
	actionName = buff.v[4]
	if isReserved(actionName) {
		return 0, ReservedError{Word: actionName}
	}

	if buff.t[5] != openBracket || buff.t[i-1] != closeBracket {
		return 0, ParseError{Expected: "[ <params...> ]",
			Found: fmt.Sprintf("%s ... %s", buff.v[5], buff.v[i-1])}
	}

	params, err := parseParams(buff.t[6:i-1], buff.v[6:i-1])
	if err != nil {
		return 0, err
	}
	action := proto.Action{
		Name:   actionName,
		From:   from,
		To:     to,
		Params: params,
	}
	b.p.Actions = append(b.p.Actions, action)
	return i, nil
}

// Parse a BSPL protocol definition from a list of tokens and values
func (b *ProtoBuilder) Parse(tokens []am.Token, values []string) error {
	i := 0
	j, err := b.parseName(tokens, values)
	if err != nil {
		return GlobalParseError{parsed: values[:i], err: err}
	}
	i += j + 1 // j+1 skip newline
	j, err = b.parseRoles(tokens[i:], values[i:])
	if err != nil {
		return GlobalParseError{parsed: values[:i], err: err}
	}
	i += j + 1
	j, err = b.parseProtoParams(tokens[i:], values[i:])
	if err != nil {
		return GlobalParseError{parsed: values[:i], err: err}
	}
	i += j + 1

	// parse first action
	j, err = b.parseActions(tokens[i:], values[i:])
	if err != nil {
		return GlobalParseError{parsed: values[:i], err: err}
	}
	i += j + 1

ACTIONS:
	for {
		if i >= len(tokens) {
			return GlobalParseError{parsed: values[:i],
				err: errors.New("Unexpected EOF")}
		}
		switch tokens[i] {
		case closeBrace:
			for k := i + 1; k < len(tokens); k++ {
				if tokens[k] != newline {
					return GlobalParseError{
						parsed: values[:i],
						err: ParseError{
							Expected: "\n",
							Found:    values[k],
						},
					}
				}
			}
			break ACTIONS
		case word:
			j, err = b.parseActions(tokens[i:], values[i:])
			if err != nil {
				return err
			}
			i += j + 1
		default:
			return ParseError{Expected: "Action or '}'", Found: values[i]}
		}
	}
	return nil
}
