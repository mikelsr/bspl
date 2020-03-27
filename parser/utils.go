package parser

import (
	"errors"
	"path/filepath"
	"runtime"

	"bitbucket.org/mikelsr/gauzaez/lexer"
	am "bitbucket.org/mikelsr/gauzaez/lexer/automaton"
)

// isReserved returns true if a string is a reserved word
func isReserved(str string) bool {
	for _, w := range reservedWords {
		if str == w {
			return true
		}
	}
	return false
}

// nextNewline returns the index of the next newline token
func nextNewline(tokens []am.Token) int {
	for i, v := range tokens {
		if string(v) == newline {
			return i
		}
	}
	return -1
}

// GetProjectDir returns the absolute path to the source code of the
// project being run
func GetProjectDir() (string, error) {
	_, fileName, _, ok := runtime.Caller(1)
	if !ok {
		return "", errors.New("Failed to locate project")
	}
	return filepath.Abs(filepath.Dir(fileName))
}

// Strip duplicates a lexer.TokenTable without whitespace tokens
// or repeated newlines. It os a costly operation but it avoids repetitive checks
func Strip(tokens lexer.TokenTable) lexer.TokenTable {
	tt := lexer.TokenTable{}
	prevToken := ""
	for i, token := range tokens.Tokens {
		c := prevToken
		prevToken = string(token)
		// skip whitespace
		if token == whitespace {
			continue
		}
		// skip duplicated newlines
		if token == newline && c == newline {
			continue
		}
		tt.Tokens = append(tt.Tokens, token)
		tt.Values = append(tt.Values, tokens.Values[i])
		tt.Lines = append(tt.Lines, tokens.Lines[i])
		tt.LinePosI = append(tt.LinePosI, tokens.LinePosI[i])
		tt.LinePosE = append(tt.LinePosE, tokens.LinePosE[i])
	}
	return tt
}
