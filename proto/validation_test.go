package proto

import (
	"testing"
)

func TestIsCircular(t *testing.T) {
	// a -> b -> c
	aa := Action{Name: "aa"}
	ab := Action{Name: "ab"}
	ac := Action{Name: "ac"}
	ad := Action{Name: "ad"}

	c := &linkedAction{action: ac}
	b := &linkedAction{action: ab, dependsOn: []*linkedAction{c}}
	a := &linkedAction{action: aa, dependsOn: []*linkedAction{b}}

	if isCircular(a) {
		t.FailNow()
	}

	// a -> b -> c -> a
	c.dependsOn = []*linkedAction{a}
	if !isCircular(a) {
		t.Fatal()
	}
	// a -> b-> c
	// \-> d -/
	c.dependsOn = nil
	d := &linkedAction{action: ad, dependsOn: []*linkedAction{c}}
	a.dependsOn = append(a.dependsOn, d)
	if isCircular(a) {
		t.Fatal()
	}
	// a -> b-> c -> a
	// \-> d -/
	c.dependsOn = []*linkedAction{a}
	if !isCircular(a) {
		t.Fatal()
	}
}

func TestCreateLinkedActions(t *testing.T) {
	p := testProtocol()
	noncircular := createLinkedActions(p.Actions)
	if circ, _ := areCircular(noncircular); circ {
		t.FailNow()
	}
	// insert circular dependency
	p.Actions[0].Params = append(p.Actions[0].Params, Parameter{Name: "price", Io: In})
	circular := createLinkedActions(p.Actions)
	if circ, _ := areCircular(circular); !circ {
		t.FailNow()
	}
}

func TestValidate(t *testing.T) {
	errMsg := "Excpected validation to fail"
	p := testProtocol()
	// correct validation
	if err := Validate(p); err != nil {
		t.Log(err)
		t.FailNow()
	}
	// insert circular dependency
	p.Actions[0].Params = append(p.Actions[0].Params, Parameter{Name: "price", Io: In})
	if err := Validate(p); err == nil {
		t.Log(errMsg)
		t.FailNow()
	}
	p.Actions = append(p.Actions, Action{
		Name:   "Fail",
		From:   p.Roles[0],
		To:     p.Roles[1],
		Params: []Parameter{{Name: "madeup", Key: true, Io: Out}},
	})
	if err := Validate(p); err == nil {
		t.Log(errMsg)
		t.FailNow()
	}
	p.Actions = []Action{{
		Name:   "Nokeyparams",
		From:   p.Roles[0],
		To:     p.Roles[1],
		Params: []Parameter{{Name: "madeup"}},
	}}
	if err := Validate(p); err == nil {
		t.Log(errMsg)
		t.FailNow()
	}
}
