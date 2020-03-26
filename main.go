package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mikelsr/bspl/parser"
)

func main() {
	dir, err := parser.GetProjectDir()
	if err != nil {
		panic(err)
	}
	bsplSource, err := os.Open(filepath.Join(dir, "test", "samples", "example_1.bspl"))
	if err != nil {
		panic(err)
	}
	tokens, err := parser.LexStream(bsplSource)
	if err != nil {
		panic(err)
	}
	fmt.Println(tokens)
	stripped := parser.Strip(*tokens)
	fmt.Printf("Stripped:\n%s", &stripped) // *TokenTable has String() method
}
