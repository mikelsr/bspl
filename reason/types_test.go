package reason

import (
	"testing"

	"github.com/mikelsr/bspl/proto"
)

func testParams() []proto.Parameter {
	return []proto.Parameter{
		{Name: "ID", Key: true, Io: proto.Out},
		{Name: "item", Io: proto.Out},
		{Name: "price", Io: proto.Out},
	}
}

func TestValues_Equal(t *testing.T) {
	params := testParams()
	strValue := "testvalue"
	v1 := make(Values)
	v2 := make(Values)
	for _, p := range params {
		v1[p.String()] = strValue
		v2[p.String()] = strValue
	}
	if !v1.Equals(v2) {
		t.Log("Couldn't compare equal values v1 and v2")
		t.FailNow()
	}

	for k := range v2 {
		v2[k] = "_"
		break
	}
	if v1.Equals(v2) {
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
	if v1.Equals(v3) {
		t.Log("Dind't fail comparing different values v1 and v3")
		t.FailNow()
	}
}
