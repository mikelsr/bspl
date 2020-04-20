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

type parameters []Parameter

func (p parameters) Len() int {
	return len(p)
}

func (p parameters) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p parameters) Less(i, j int) bool {
	p1, p2 := p[i], p[j]
	// i is key, j is not
	if p1.Key {
		if !p2.Key {
			return true
		}
	}
	// j is key, i is not
	if p2.Key {
		if !p1.Key {
			return false
		}
	}
	// If IO is different, check order:
	// 1. In, 2. Nil, 3. Out
	if p1.Io != p2.Io {
		switch p1.Io {
		case In:
			return true
		case Nil:
			return p2.Io == Out
		default:
		case Out:
			return false
		}
	}
	// If IO is the same, sort alphabetically
	return p1.Name < p2.Name
}

// SortParameters sorts protocol parameters
// 1 - Keys sorted alphabetically
// 2 - Ins sorted alphabetically
// 3 - Nils sorted alphabetically
// 4 - Outs sorted alphabetically
func SortParameters(params []Parameter) {
	sort.Sort(parameters(params))
}

// SortActions sorts actions alphabetically
func SortActions(acts []Action) {
	sort.Sort(actions(acts))
}

// SortRoles sorts roles alphabetically
func SortRoles(rols []Role) {
	sort.Sort(roles(rols))
}
