package implementation

import (
	"bytes"
	"testing"

	"github.com/mikelsr/bspl/proto"
)

func TestMarshalAction(t *testing.T) {
	expected := []byte{66, 117, 121, 101, 114, 32, 45, 62, 32, 83, 101, 108,
		108, 101, 114, 58, 32, 79, 102, 102, 101, 114, 91, 105, 110, 32,
		73, 68, 32, 107, 101, 121, 44, 32, 105, 110, 32, 105, 116, 101,
		109, 44, 32, 111, 117, 116, 32, 112, 114, 105, 99, 101, 93}
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

func TestInstance_MarshalAndUnmarshal(t *testing.T) {
	testInstanceMarshal(t)
	testInstanceUnmarshal(t)
}

func testInstanceMarshal(t *testing.T) {
	i := testInstance()
	_, err := i.Marshal()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
}

func testInstanceUnmarshal(t *testing.T) {
	expected := testInstance()
	data, _ := expected.Marshal()
	var i Instance
	var err error
	if i, err = i.Unmarshal(data); err != nil {
		t.Log(err)
		t.FailNow()
	}
	if !expected.Equals(i) {
		t.FailNow()
	}
}

func TestMessage_MarshalAndUnmarshal(t *testing.T) {
	testMessageMarshal(t)
	testInstanceUnmarshal(t)
}

func testMessageMarshal(t *testing.T) {
	expected := []byte{123, 34, 105, 110, 115, 116, 97, 110, 99, 101, 95, 107,
		101, 121, 34, 58, 34, 80, 114, 111, 116, 111, 78, 97, 109, 101, 44,
		73, 68, 58, 88, 34, 44, 34, 97, 99, 116, 105, 111, 110, 34, 58, 34,
		66, 117, 121, 101, 114, 32, 45, 92, 117, 48, 48, 51, 101, 32, 83,
		101, 108, 108, 101, 114, 58, 32, 79, 102, 102, 101, 114, 91, 105,
		110, 32, 73, 68, 32, 107, 101, 121, 44, 32, 105, 110, 32, 105, 116,
		101, 109, 44, 32, 111, 117, 116, 32, 112, 114, 105, 99, 101, 93, 34,
		44, 34, 118, 97, 108, 117, 101, 115, 34, 58, 123, 34, 105, 110, 32,
		73, 68, 32, 107, 101, 121, 34, 58, 34, 88, 34, 44, 34, 105, 110, 32,
		105, 116, 101, 109, 34, 58, 34, 88, 34, 44, 34, 111, 117, 116, 32,
		112, 114, 105, 99, 101, 34, 58, 34, 88, 34, 125, 125}
	i := testInstance()
	var m Message
	for _, v := range i.Messages() {
		m = v.(Message)
		break
	}
	data, err := m.Marshal()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if !bytes.Equal(data, expected) {
		t.FailNow()
	}
}

func testMessageUnarshal(t *testing.T) {
	i := testInstance()
	var expected Message
	for _, v := range i.Messages() {
		expected = v.(Message)
		break
	}
	data, _ := expected.Marshal()
	var m Message
	im, err := m.Unmarshal(data)
	if err != nil {
		t.FailNow()
	}
	m = im.(Message)
	if !compareMessages(
		Messages{m.Action().String(): m},
		Messages{expected.Action().String(): expected}) {
		t.FailNow()
	}
}
