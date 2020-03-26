package parser

import (
	"errors"
	"path/filepath"
	"runtime"

	"bitbucket.org/mikelsr/gauzaez/lexer"
)

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
// it's a costly operation but it avoids repetitive checks
func Strip(tokens lexer.TokenTable) lexer.TokenTable {
	tt := lexer.TokenTable{}
	for i, token := range tokens.Tokens {
		if token != "whitespace" {
			tt.Tokens = append(tt.Tokens, token)
			tt.Values = append(tt.Values, tokens.Values[i])
			tt.Lines = append(tt.Lines, tokens.Lines[i])
			tt.LinePosI = append(tt.LinePosI, tokens.LinePosI[i])
			tt.LinePosE = append(tt.LinePosE, tokens.LinePosE[i])
		}
	}
	return tt
}
