package parser

import (
	"io"
	"path/filepath"

	"bitbucket.org/mikelsr/gauzaez/lexer"
)

// newLexer creates a lexer.Lexer with the lexer.Rules taken from the
// config/lexer.json file found in this project
func newLexer() (*lexer.Lexer, error) {
	dir, err := GetProjectDir()
	path := filepath.SplitList(dir)
	dir = filepath.Join(path[0 : len(path)-1]...)
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
