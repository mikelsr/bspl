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

func TestProtocol_Dependencies(t *testing.T) {
	p := testProtocol()
	request, offer := p.Actions[0], p.Actions[1]
	deps := p.Dependencies(offer)
	if len(deps) != 1 {
		t.FailNow()
	}
	if deps[0].String() != request.String() {
		t.FailNow()
	}
}
