package bspl

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mikelsr/bspl/parser"
)

func TestParse(t *testing.T) {
	dir, err := parser.GetProjectDir()
	if err != nil {
		panic(err)
	}
	bsplSource, err := os.Open(filepath.Join(dir, "test", "samples", "example_1.bspl"))
	if err != nil {
		panic(err)
	}
	if _, err := Parse(bsplSource); err != nil {
		t.FailNow()
	}
}
