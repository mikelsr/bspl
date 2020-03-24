package parser

import (
	"io"
	"path/filepath"

	"bitbucket.org/mikelsr/gauzaez/lexer"
)

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

func LexStream(in io.Reader) (*lexer.TokenTable, error) {
	lex, err := newLexer()
	if err != nil {
		return nil, err
	}
	return lex.Tokenize(in)
}
