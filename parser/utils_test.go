package parser

import (
	"testing"

	am "bitbucket.org/mikelsr/gauzaez/lexer/automaton"
)

func TestNextNewLine(t *testing.T) {
	tokens := []am.Token{word, word, newline}
	if nextNewline(tokens) != 2 {
		t.FailNow()
	}
	tokens[0] = newline
	if nextNewline(tokens) != 0 {
		t.FailNow()
	}
	tokens = []am.Token{arrow, colon, comma}
	if nextNewline(tokens) != -1 {
		t.FailNow()
	}
}

func TestStrip(t *testing.T) {
	count := 0
	prevToken := ""
	for _, token := range testTokensW.Tokens {
		c := prevToken
		prevToken = string(token)
		if token == whitespace {
			count++
		}
		if token == newline && c == newline {
			count++
		}
	}
	strippedTokens := Strip(*testTokensW)
	if len(testTokensW.Tokens) != len(strippedTokens.Tokens)+count {
		t.FailNow()
	}
	for _, token := range strippedTokens.Tokens {
		if token == whitespace {
			t.FailNow()
		}
	}
}
