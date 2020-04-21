package instance

import (
	"bytes"
	"testing"

	"github.com/mikelsr/bspl/proto"
)

func TestMarshalAction(t *testing.T) {
	expected := []byte{66, 117, 121, 101, 114, 32, 45, 62, 32, 83, 101, 108, 108, 101,
		114, 58, 32, 82, 101, 113, 117, 101, 115, 116, 91, 111, 117, 116, 32, 73,
		68, 32, 107, 101, 121, 44, 32, 111, 117, 116, 32, 105, 116, 101, 109, 93}
	sample := testProtocol().Actions[0]
	marshalled := MarshalAction(sample)
	if !bytes.Equal(marshalled, expected) {
		t.FailNow()
	}
}

func TestUnmarshalAction(t *testing.T) {
	sample := testProtocol().Actions[0]
	action := &proto.Action{}
	if err := UnmarshalAction(action, []byte(sample.String())); err != nil {
		t.Log(err)
		t.FailNow()
	}
	if action.String() != sample.String() {
		t.FailNow()
	}
	if err := UnmarshalAction(action, []byte(sample.String())[:5]); err == nil {
		t.FailNow()
	}
}
