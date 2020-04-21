package bspl

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mikelsr/bspl/parser"
)

func TestCompare(t *testing.T) {
	dir, err := parser.GetProjectDir()
	if err != nil {
		panic(err)
	}
	bsplSourceA, err := os.Open(filepath.Join(dir, "test", "samples", "example_1.bspl"))
	if err != nil {
		panic(err)
	}
	bsplSourceB, err := os.Open(filepath.Join(dir, "test", "samples", "example_1.bspl"))
	if err != nil {
		panic(err)
	}
	protoA, errA := Parse(bsplSourceA)
	protoB, errB := Parse(bsplSourceB)
	if errA != nil || errB != nil {
		t.FailNow()
	}
	if !Compare(protoA, protoB) {
		t.FailNow()
	}
	protoB.Roles[0] = Role(protoB.Roles[0] + "_")
	protoB.Actions[0], protoB.Actions[1] = protoB.Actions[1], protoB.Actions[0]
	if Compare(protoA, protoB) {
		t.FailNow()
	}
}
