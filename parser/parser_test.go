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

func TestParse(t *testing.T) {
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
	if _, err := Parse(bsplSource); err != nil {
		t.FailNow()
	}
}

func TestProtoBuilder_parseName(t *testing.T) {
	b := new(ProtoBuilder)
	i, err := b.parseName(testTokens.Tokens, testTokens.Values)
	if err != nil || i != 2 {
		t.FailNow()
	}
	errTokens := []am.Token{newline}
	errValues := []string{"\n"}
	if i, err = b.parseName(errTokens, errValues); err == nil || i != 0 {
		t.FailNow()
	}
	// invalid syntax
	errTokens = []am.Token{openBracket, word, newline}
	errValues = []string{"{", "name", "\n"}
	if i, err = b.parseName(errTokens, errValues); err == nil || i != 0 {
		t.FailNow()
	}
	// use a reserved word as a name
	errTokens = []am.Token{word, openBracket, newline}
	errValues = []string{"key", "{", "\n"}
	if i, err = b.parseName(errTokens, errValues); err == nil || i != 0 {
		t.FailNow()
	}
}

func TestProtoBuilder_parseRole(t *testing.T) {
	b := new(ProtoBuilder)
	tokens := []am.Token{word, word, comma, word, comma, word, newline}
	values := []string{Role, "A", ",", "B", ",", "C", "\n"}
	if i, err := b.parseRoles(tokens, values); err != nil || i != len(tokens)-1 {
		t.FailNow()
	}
	if !reflect.DeepEqual(b.p.Roles, []proto.Role{"A", "B", "C"}) {
		t.FailNow()
	}
	// one of the roles is a reserved keyword
	values[3] = Role
	if _, err := b.parseRoles(tokens, values); err == nil {
		t.FailNow()
	}
	// invalid syntax
	tokens = []am.Token{word, word, comma, word, comma, newline}
	values = []string{"role", "A", ",", "B", ",", "\n"}
	if _, err := b.parseRoles(tokens, values); err == nil {
		t.FailNow()
	}
	// missing comma
	tokens = []am.Token{word, word, word, newline}
	values = []string{"role", "A", "B", "C", "\n"}
	if _, err := b.parseRoles(tokens, values); err == nil {
		t.FailNow()
	}
	// missing role reserved word
	tokens = []am.Token{word, word, comma, word, newline}
	values = []string{"ERR", "A", ",", "B", "\n"}
	if _, err := b.parseRoles(tokens, values); err == nil {
		t.FailNow()
	}
	// repeated role
	tokens = []am.Token{word, word, comma, word, newline}
	values = []string{"role", "A", ",", "A", "\n"}
	if _, err := b.parseRoles(tokens, values); err == nil {
		t.FailNow()
	}
}

func TestParseParams(t *testing.T) {
	tokens := []am.Token{word, word, word, comma, word, word, comma, word, word, comma, word}
	values := []string{In, "a", Key, ",", Out, "b", ",", "c", Key, ",", "d"}
	expected := []proto.Parameter{
		{Io: proto.IO(In), Name: "a", Key: true},
		{Io: proto.IO(Out), Name: "b", Key: false},
		{Io: proto.IO(Nil), Name: "c", Key: true},
		{Io: proto.IO(Nil), Name: "d", Key: false},
	}
	params, err := parseParams(tokens, values)
	if err != nil || !reflect.DeepEqual(params, expected) {
		t.FailNow()
	}
	values[0] = Key
	if _, err = parseParams(tokens, values); err == nil {
		t.FailNow()
	}
	tokens[0] = arrow
	if _, err = parseParams(tokens, values); err == nil {
		t.FailNow()
	}
	tokens[0] = word
	values[0] = In
	// multiple params without comma separation
	tokens[3] = word
	if _, err = parseParams(tokens, values); err == nil {
		t.FailNow()
	}
	tokens[3] = comma
	// repeat parameter name
	values[5] = "a"
	if _, err = parseParams(tokens, values); err == nil {
		t.FailNow()
	}
	values[5] = "b"
	// assign a reserved keyword to param names from last to first so all
	// are checked
	for _, i := range []int{10, 7, 5, 1} {
		values[i] = Role // reserved keyword
		if _, err = parseParams(tokens, values); err == nil {
			t.FailNow()
		}
	}
}

func TestProtoBuilder_parseProtoParams(t *testing.T) {
	b := new(ProtoBuilder)
	tokens := []am.Token{word, word, word, word, comma, word, word, comma, word, word, newline}
	values := []string{Param, Out, "ID", Key, ",", Out, "out_param", ",", In, "in_param", "\n"}
	if i, err := b.parseProtoParams(tokens, values); err != nil || i != len(tokens)-1 {
		t.FailNow()
	}
	expected := []proto.Parameter{
		{Io: proto.IO(Out), Name: "ID", Key: true},
		{Io: proto.IO(Out), Name: "out_param", Key: false},
		{Io: proto.IO(In), Name: "in_param", Key: false},
	}
	if !reflect.DeepEqual(b.p.Parameters(), expected) {
		t.FailNow()
	}
	values[0] = Role
	if _, err := b.parseProtoParams(tokens, values); err == nil {
		t.FailNow()
	}
	values[0] = Param
	tokens[0] = comma
	if _, err := b.parseProtoParams(tokens, values); err == nil {
		t.FailNow()
	}
	tokens[0] = word
	tokens[3] = comma
	if _, err := b.parseProtoParams(tokens, values); err == nil {
		t.FailNow()
	}
}

func TestProtoBuilder_parseActions(t *testing.T) {
	b := new(ProtoBuilder)
	b.p.Roles = []proto.Role{"From", "To"}

	tokens := []am.Token{word, arrow, word, colon, word, openBracket, word, comma, word, word, closeBracket, newline}
	values := []string{"From", "->", "To", ":", "Action", "[", "ID", ",", In, "in_param", "]", "\n"}
	if i, err := b.parseActions(tokens, values); err != nil || i != len(tokens)-1 {
		t.FailNow()
	}
}

func TestProtoBuilder_Parse(t *testing.T) {
	b := new(ProtoBuilder)
	if err := b.Parse(testTokens.Tokens, testTokens.Values); err != nil {
		t.Log(err)
		t.FailNow()
	}
}
