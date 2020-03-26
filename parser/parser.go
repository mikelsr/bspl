package parser

import (
	"bitbucket.org/mikelsr/gauzaez/lexer"
	am "bitbucket.org/mikelsr/gauzaez/lexer/automaton"
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
)

var (
	// reservedWords is a list of words reserved by the parser
	reservedWords = []string{In, Key, Nil, Out, Param}
	// nameStructure is the token structure expected for name parsing
	nameStructure = []expected{
		eToken{name: word},        // protocol name
		eToken{name: openBracket}, // start protocol definition
	}
	// roleStructure is the token structure expected for role parsing
	roleStructure = []expected{
		eGroup{
			expected: []expected{
				eToken{name: word},
				eToken{name: comma},
			},
			optional: true, // in case of single role // TODO: is this possible?
			repeats:  true, // in case of >2 roles, this group is repeated
		},
		eToken{name: word}, // final role
	}
	// paramStructure is the token structure expected for parameter parsing
	paramStructure = []expected{
		eGroup{
			expected: []expected{
				eToken{name: word, optional: true, reserved: true, mustBe: []string{In, Nil, Out}}, // in, out, nil
				eToken{name: word}, // value
				eToken{name: word, optional: true, reserved: true, mustBe: []string{Key}}, // key
			}, repeats: true,
		},
	}
	// protoParamStructure is the structure expected at the protocol
	// parameter declaration, no the action parameter declaration
	protoParamStructure = append([]expected{
		eToken{name: word, reserved: true, mustBe: []string{Param}}, // parameter
	}, paramStructure...) // list of parameters
	// actionStructure is the token structure expected for action parsing
	actionStructure = append([]expected{
		eToken{name: word},                      // from
		eToken{name: arrow}, eToken{name: word}, // -> to
		eToken{name: colon}, eToken{name: word}, // : actionName
		eToken{name: openBrace}, // [ (open parameters)
	}, append(
		paramStructure,           // parameters
		eToken{name: closeBrace}, // ] (close brace)
	)...)
)

type expected interface {
	isAtom() bool
	isOptional() bool
	mayRepeat() bool
}

type eToken struct {
	name     string
	optional bool
	repeats  bool
	reserved bool
	mustBe   []string // if the token is a word it MUST be one of this values
}

func (e eToken) isAtom() bool {
	return true
}

func (e eToken) isOptional() bool {
	return e.optional
}

func (e eToken) mayRepeat() bool {
	return e.repeats
}

type eGroup struct {
	expected []expected
	optional bool
	repeats  bool
}

func (eg eGroup) isAtom() bool {
	return false
}

func (eg eGroup) isOptional() bool {
	return eg.optional
}

func (eg eGroup) mayRepeat() bool {
	return eg.repeats
}

// validateToken checks if a token is of the expected type,
// is a keyword with an appropiate value or is a non-keyword
// that doesn't take a reserved value
func validateToken(token string, value string, e eToken) bool {
	if token != e.name {
		return false
	}
	// reserved word
	if e.reserved {
		// check that it has a reserved value
		ok := false
		for _, v := range reservedWords {
			if v == value {
				ok = true
			}
		}
		if !ok {
			return false
		}
		// not a reserved word
	} else {
		// check that it doesn't have the value of a reserved word
		for _, v := range reservedWords {
			if v == value {
				return false
			}
		}
	}
	// check that it has the expected value
	if len(e.mustBe) > 0 {
		ok := false
		for _, v := range e.mustBe {
			if v == value {
				ok = true
			}
		}
		if !ok {
			return false
		}
	}
	return true
}

// optionalForward recursively checks if optional tokens are being skipped
func optionalForward(token, value string, es []expected) (bool, int) {
	// if it is a group, validation failed
	if !es[0].isAtom() {
		return false, 1
	}
	e := es[0].(eToken)
	// validate current token
	if validateToken(token, value, e) {
		return true, 1
	}
	// if optional, validate next token
	if e.isOptional() {
		ok, i := optionalForward(token, value, es[1:])
		return ok, 1 + i
	}
	// if not optional and not validated, token is invalid
	return false, 1
}

// validateGroup checks a tokens against a token group (eGroup)
func validateGroup(tokens []am.Token, values []string, group eGroup) (bool, int) {
	i := 0
	for i < len(group.expected) {
		next := group.expected[i]
		if next.isAtom() {
			e := next.(eToken)
			token := string(tokens[i])
			value := values[i]
			if !validateToken(token, value, e) {
				if e.isOptional() {
					// recursively check ahead if optional tokens are being skipped
					ok, j := optionalForward(token, value, group.expected[i:])
					if !ok {
						// check backwards
						return false, i + j
					}
					i += j
				}
			}
		} else {
			g := next.(eGroup)
			ok, j := validateGroup(tokens[i:], values[i:], g)
			i = i + j
			if !ok {
				return false, i
			}
		}
		i++ // TODO
	}
	return true, i
}

// parseSection parses a token table to match an expected structure until
// finding a separator. It starts at the lexer.TokenTable index:index
func validateSection(tokens lexer.TokenTable, structure []expected, separator string, index int) (bool, int) {
	current := struct{ t, v string }{
		t: string(tokens.Tokens[index]),
		v: tokens.Values[index],
	}
	e := structure[0]
	for current.t != separator {
		// validate a single expected token
		if e.isAtom() {
			if !validateToken(current.t, current.v, e.(eToken)) {
				return false, index
			}
		}
	}
	return true, index
}
