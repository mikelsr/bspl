package proto

import (
	"errors"
	"fmt"
)

// Validate a Protocol
func Validate(p Protocol) error {
	keyParams := p.Keys()
	if len(keyParams) < 1 {
		return errors.New("No key parameters")
	}

	// check that actions have at least one key and that key
	// has been declared in the protocol parameters
	for _, a := range p.Actions {
		// check that both roles are defined in the protocol
		for _, actionRole := range []Role{a.From, a.To} {
			definedRole := false
			for _, role := range p.Roles {
				if role == actionRole {
					definedRole = true
				}
			}
			if !definedRole {
				return ValidationError{Err: fmt.Errorf("Unknown role: %s", actionRole)}
			}
		}
		found := false
		// check it at least one action parameter is a key protocol parameter
	KeyCheck:
		for _, k := range a.Params {
			for _, pk := range keyParams {
				if pk.Name == k.Name {
					found = true
					break KeyCheck
				}
			}
		}
		if !found {
			return ValidationError{Err: fmt.Errorf(
				"Action '%s' has no key parameters in common with '%s'",
				a.Name, p.Name)}
		}
	}
	dependencies := createLinkedActions(p.Actions)
	if circ, a := areCircular(dependencies); circ {
		return ValidationError{Err: fmt.Errorf(
			"Circular dependency at action '%s'", a.Name)}
	}
	return nil
}

type linkedAction struct {
	action    Action
	dependsOn []*linkedAction
}

func intersection(ins []Parameter, outs []Parameter) []Parameter {
	inter := make([]Parameter, 0)
	for _, out := range outs {
		for _, in := range ins {
			if in.Name == out.Name {
				inter = append(inter, out)
			}
		}
	}
	return inter
}

func findLinkedAction(la []*linkedAction, a Action) *linkedAction {
	for _, l := range la {
		if l.action.Name == a.Name {
			return l
		}
	}
	return nil
}

// this has _a lot_ of room for improvement
func createLinkedActions(as []Action) []*linkedAction {
	dependencies := make([]*linkedAction, 0)
	// each actions looks for the actions that preceed them
	for i, a := range as {
		for j, b := range as {
			if i == j {
				continue
			}
			// find common parameters between A inputs and B outputs
			inter := intersection(a.Ins(), b.Outs())
			if len(inter) > 0 {
				var dependsOn, src *linkedAction
				// find link to A if it was already created
				dependsOn = findLinkedAction(dependencies, b)
				// if not, create link to B
				if dependsOn == nil {
					dependsOn = &linkedAction{action: b}
					dependencies = append(dependencies, dependsOn)
				}
				// find link to B if it was already created
				src = findLinkedAction(dependencies, a)
				// if not found, create link to A
				if src == nil {
					src = &linkedAction{action: a, dependsOn: []*linkedAction{dependsOn}}
					dependencies = append(dependencies, src)
					// if found, modify link to A to add the new destination
				} else {
					src.dependsOn = append(src.dependsOn, dependsOn)
				}
			}
		}
	}
	return dependencies
}

func isChecked(l *linkedAction, checked []*linkedAction) bool {
	for _, x := range checked {
		if x == l {
			return true
		}
	}
	return false
}

func areCircular(la []*linkedAction) (bool, Action) {
	// for each node, follow each link until the end or
	// returning to self
	// checking each node more than once is not effective but
	// sad things happen and we must move on
	checked := make([]*linkedAction, 0)
	for _, d := range la {
		// d was already checked when calling a previous isCircular()
		if isChecked(d, checked) {
			continue
		}
		if isCircular(d) {
			return true, d.action
		}
	}
	return false, Action{}
}

func isCircular(l *linkedAction) bool {
	return _isCircular(l, []*linkedAction{})
}

func _isCircular(l *linkedAction, checked []*linkedAction) bool {
	if l == nil {
		panic("Invalid nil linkedAction")
	}
	if isChecked(l, checked) {
		return true
	}
	checked = append(checked, l)
	if len(l.dependsOn) == 0 {
		return false
	}

	// for each next node, repeat
	for _, dependsOn := range l.dependsOn {
		if _isCircular(dependsOn, checked) {
			return true
		}
	}
	return false
}
