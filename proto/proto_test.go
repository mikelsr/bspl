package proto

import (
	"testing"
)

func TestProtocol_Key(t *testing.T) {
	expected := "ProtoName,ID"
	p := testProtocol()
	if p.Key() != expected {
		t.FailNow()
	}
}
