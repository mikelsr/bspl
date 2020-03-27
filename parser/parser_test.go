package parser

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"bitbucket.org/mikelsr/gauzaez/lexer"
	am "bitbucket.org/mikelsr/gauzaez/lexer/automaton"
	"github.com/mikelsr/bspl/proto"
)

var (
	// test tokens with whitespaces
	testTokensW *lexer.TokenTable
	// test tokens without whitespaces
	testTokens *lexer.TokenTable
)

func TestMain(m *testing.M) {
	dir, err := GetProjectDir()
	if err != nil {
		panic(err)
	}
	path := strings.Split(dir, string(os.PathSeparator))
	dir = "/" + filepath.Join(path[:len(path)-1]...)
	dir = filepath.Join(dir, "test", "samples", "example_1.bspl")
	bsplSource, err := os.Open(dir)
	if err != nil {
		panic(err)
	}
	testTokensW, err = LexStream(bsplSource)
	if err != nil {
		panic(err)
	}
	// Strip token table once to speed up tests
	tt := Strip(*testTokensW)
	testTokens = &tt
	m.Run()
}

func TestProtoBuilder_parseName(t *testing.T) {
	b := new(ProtoBuilder)
	i, ok := b.parseName(testTokens.Tokens, testTokens.Values)
	if !ok || i != 2 {
		t.FailNow()
	}
	errTokens := []am.Token{newline}
	errValues := []string{"\n"}
	if i, ok = b.parseName(errTokens, errValues); ok || i != 0 {
		t.FailNow()
	}
	// invalid syntax
	errTokens = []am.Token{openBracket, word, newline}
	errValues = []string{"{", "name", "\n"}
	if i, ok = b.parseName(errTokens, errValues); ok || i != 0 {
		t.FailNow()
	}
	// use a reserved word as a name
	errTokens = []am.Token{word, openBracket, newline}
	errValues = []string{"key", "{", "\n"}
	if i, ok = b.parseName(errTokens, errValues); ok || i != 0 {
		t.FailNow()
	}
}

func TestProtoBuilder_parseRole(t *testing.T) {
	b := new(ProtoBuilder)
	tokens := []am.Token{word, word, comma, word, comma, word, newline}
	values := []string{"role", "A", ",", "B", ",", "C", "\n"}
	if i, ok := b.parseRoles(tokens, values); !ok || i != len(tokens)-1 {
		t.FailNow()
	}
	if !reflect.DeepEqual(b.p.Roles, []proto.Role{"A", "B", "C"}) {
		t.FailNow()
	}
	// one of the roles is a reserved keyword
	values[3] = Role
	if _, ok := b.parseRoles(tokens, values); ok {
		t.FailNow()
	}
	// invalid syntax
	tokens = []am.Token{word, word, comma, word, comma, newline}
	values = []string{"role", "A", ",", "B", ",", "\n"}
	if _, ok := b.parseRoles(tokens, values); ok {
		t.FailNow()
	}
	// missing comma
	tokens = []am.Token{word, word, word, newline}
	values = []string{"role", "A", "B", "C", "\n"}
	if _, ok := b.parseRoles(tokens, values); ok {
		t.FailNow()
	}
	// missing role reserved word
	tokens = []am.Token{word, word, comma, word, newline}
	values = []string{"ERR", "A", ",", "B", "\n"}
	if _, ok := b.parseRoles(tokens, values); ok {
		t.FailNow()
	}
	// repeated role
	tokens = []am.Token{word, word, comma, word, newline}
	values = []string{"role", "A", ",", "A", "\n"}
	if _, ok := b.parseRoles(tokens, values); ok {
		t.FailNow()
	}
}
