package proto

import (
	"sort"
)

// actions implements (Go sort.Interface) to sort actions alphabetically
type actions []Action

func (a actions) Len() int {
	return len(a)
}

func (a actions) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a actions) Less(i, j int) bool {
	return a[i].String() < a[j].String()
}

// roles implements (Go sort.Interface) to sort roles alphabetically
type roles []Role

func (r roles) Len() int {
	return len(r)
}

func (r roles) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r roles) Less(i, j int) bool {
	return r[i] < r[j]
}

// SortActions sorts actions alphabetically
func (p *Protocol) SortActions() {
	sort.Sort(actions(p.Actions))
}

// SortRoles sorts roles alphabetically
func (p *Protocol) SortRoles() {
	sort.Sort(roles(p.Roles))
}
