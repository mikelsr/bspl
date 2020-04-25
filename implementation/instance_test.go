package implementation

import (
	"testing"

	"github.com/mikelsr/bspl/proto"
)

func TestIntance_Key(t *testing.T) {
	expected := "ProtoName,ID:X"
	key := testInstance().Key()
	if expected != key {
		t.FailNow()
	}
}

func TestInstance_Equals(t *testing.T) {
	i1 := testInstance()
	i2 := testInstance()
	if !i1.Equals(i2) {
		t.FailNow()
	}
}

func TestInstance_Diff(t *testing.T) {
	p := testProtocol()
	roles := Roles{
		proto.Role("Buyer"):  "B",
		proto.Role("Seller"): "S",
	}
	// i1 is the empty instance
	i1 := NewInstance(p, roles)
	// i2 is the same as i1 but after running "Request"
	i2 := NewInstance(p, roles)
	i2.SetValue("ID", "testID")
	i2.SetValue("item", "testItem")
	actions, diff, err := i1.Diff(i2)
	if err != nil {
		t.Error(err)
	}
	if len(actions) != 1 {
		t.Fatal("Missing actions")
	}
	if actions[0].Name != "Request" {
		t.Fatal("Wrong action name")
	}
	if len(diff) != 2 {
		t.Fatal("Missing parameters")
	}
}

func TestInstace_Update(t *testing.T) {
	p := testProtocol()
	roles := Roles{
		proto.Role("Buyer"):  "B",
		proto.Role("Seller"): "S",
	}
	// i1 is the empty instance
	i1 := NewInstance(p, roles)
	// i2 is the same as i1 but after running "Request"
	i2 := NewInstance(p, roles)
	i2.SetValue("ID", "testID")
	i2.SetValue("item", "testItem")
	// i3 is the same as i2 but after running "Offer"
	i3 := NewInstance(p, roles)
	i3.SetValue("price", "testPrice")

	if err := i1.Update(i2); err != nil {
		t.Error(err)
	}
	if !i1.Equals(i2) {
		t.Fatal("i1 and i2 differ after update")
	}
}
