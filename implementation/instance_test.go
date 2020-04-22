package implementation

import (
	"testing"
)

func TestIntance_Key(t *testing.T) {
	expected := "ProtoName,ID:X"
	key := testInstance().Key()
	if expected != key {
		t.FailNow()
	}
}

func TestInstance_Equales(t *testing.T) {
	i1 := testInstance()
	i2 := testInstance()
	if !i1.Equals(i2) {
		t.FailNow()
	}
}
