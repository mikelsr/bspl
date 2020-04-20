package proto

import (
	"testing"
)

func TestSortParameters(t *testing.T) {
	params := []Parameter{
		{Name: "E", Io: Nil},
		{Name: "D", Io: In},
		{Name: "F", Io: Out},
		{Name: "B", Key: true, Io: Nil},
		{Name: "A", Key: true, Io: In},
		{Name: "C", Key: true, Io: Out},
	}
	expected := []string{"A", "B", "C", "D", "E", "F"}
	SortParameters(params)
	for i, p := range params {
		if p.Name != expected[i] {
			t.FailNow()
		}
	}
}

func TestSortActions(t *testing.T) {
	acts := []Action{
		{Name: "B", From: "B", To: "A", Params: []Parameter{
			{Name: "ID", Key: true, Io: Out},
		}},
		{Name: "A", From: "B", To: "A", Params: []Parameter{
			{Name: "ID", Key: true, Io: Out},
		}},
	}
	expected := []string{"A", "B"}
	for i, a := range acts {
		if a.Name != expected[i] {
			t.FailNow()
		}
	}
}

func TestSortRoles(t *testing.T) {
	r := []Role{"B", "A"}
	expected := []Role{"A", "B"}
	SortRoles(r)
	for i, role := range r {
		if role != expected[i] {
			t.FailNow()
		}
	}
}
