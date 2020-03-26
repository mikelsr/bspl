package parser

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"bitbucket.org/mikelsr/gauzaez/lexer"
)

// newLexer creates a lexer.Lexer with the lexer.Rules taken from the
// config/lexer.json file found in this project
func newLexer() (*lexer.Lexer, error) {
	dir, err := GetProjectDir()
	path := strings.Split(dir, string(os.PathSeparator))
	dir = "/" + filepath.Join(path[:len(path)-1]...)
	if err != nil {
		return nil, err
	}
	// use the rules in the config folder of the project
	rules, err := lexer.MakeRules(filepath.Join(dir, "config", "lexer.json"))
	if err != nil {
		return nil, err
	}
	return lexer.MakeLexer(*rules)
}

// LexStream passes the input stream throgh an anonimous lexer
func LexStream(in io.Reader) (*lexer.TokenTable, error) {
	lex, err := newLexer()
	if err != nil {
		return nil, err
	}
	return lex.Tokenize(in)
}
