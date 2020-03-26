package parser

import "testing"

func TestStrip(t *testing.T) {
	count := 0
	for _, token := range testTokens.Tokens {
		if token == whitespace {
			count++
		}
	}
	strippedTokens := Strip(*testTokens)
	if len(testTokens.Tokens) != len(strippedTokens.Tokens)+count {
		t.FailNow()
	}
	for _, token := range strippedTokens.Tokens {
		if token == whitespace {
			t.FailNow()
		}
	}
}
