package reason

import (
	"github.com/mikelsr/bspl/proto"
)

// Values maps component (Protocol, Action...) identifiers to
// string values
type Values map[string]string

// Roles maps protocol roles to string IDs of the agents
// assuming the roles
type Roles map[proto.Role]string

// Equals compares two Values (Values = map[string]string)
func (v Values) Equals(others Values) bool {
	if len(v) != len(others) {
		return false
	}
	for k, v1 := range v {
		if v2, found := others[k]; !found || v1 != v2 {
			return false
		}
	}
	return true
}
