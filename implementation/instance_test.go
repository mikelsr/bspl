package implementation

import (
	"testing"
)

func TestCompareValues(t *testing.T) {
	p := testProtocol()
	strValue := "testvalue"
	v1 := make(Values)
	v2 := make(Values)
	for _, p := range p.Parameters() {
		v1[p.String()] = strValue
		v2[p.String()] = strValue
	}
	if !compareValues(v1, v2) {
		t.Log("Couldn't compare equal values v1 and v2")
		t.FailNow()
	}

	for k := range v2 {
		v2[k] = "_"
		break
	}
	if compareValues(v1, v2) {
		t.Log("Dind't fail comparing different values v1 and v2")
		t.FailNow()
	}
	v3 := make(Values)
	n := len(v1)
	i := 0
	for k, v := range v1 {
		if i == n-1 {
			break
		}
		i++
		v3[k] = v
	}
	if compareValues(v1, v3) {
		t.Log("Dind't fail comparing different values v1 and v3")
		t.FailNow()
	}
}

func TestCompareMessages(t *testing.T) {

}

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
