package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"bitbucket.org/mikelsr/gauzaez/lexer"
	am "bitbucket.org/mikelsr/gauzaez/lexer/automaton"
)

var testTokens *lexer.TokenTable

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
	testTokens, err = LexStream(bsplSource)
	if err != nil {
		panic(err)
	}
	m.Run()
}

func TestOtionalForward(t *testing.T) {
	ex := []expected{
		eToken{name: word, optional: true, reserved: true, mustBe: []string{In, Nil, Out}}, // in, out, nil
		eToken{name: word}, // value
		eToken{name: word, optional: true, reserved: true, mustBe: []string{Key}}, // key
	}

	tt := &lexer.TokenTable{}
	// skip scope in parameter declaration, use only name and key
	tt.Tokens = []am.Token{word, word}
	tt.Values = []string{"param_name", Key}
	ok, i := optionalForward(string(tt.Tokens[0]), tt.Values[0], ex)
	if !ok || i != 2 {
		t.FailNow()
	}
	// skip parameter name, should fail validation
	tt.Values = []string{Nil, Key} // missing compulsory world
	if ok, _ := optionalForward(string(tt.Tokens[1]), tt.Values[1], ex[1:]); ok {
		t.FailNow()
	}
}

func TestValidateToken(t *testing.T) {
	et := eToken{name: word, reserved: true}
	if !validateToken(word, Key, et) {
		t.FailNow()
	}
	if validateToken(word, "_", et) {
		t.FailNow()
	}
	et.mustBe = []string{In, Out, Nil}
	if validateToken(word, Key, et) || !validateToken(word, In, et) {
		t.FailNow()
	}
	et.reserved = false
	if validateToken(word, Key, et) {
		t.FailNow()
	}
	et.mustBe = []string{"A", "B", "C"}
	if !validateToken(word, "B", et) {
		t.FailNow()
	}
	et = eToken{name: colon}
	if !validateToken(colon, ":", et) || validateToken(word, ":", et) {
		t.FailNow()
	}
}
