package implementation

import (
	"reflect"

	"github.com/mikelsr/bspl/reason"
)

// Messages maps Instance.Key() to Message
type Messages = reason.Messages

// Roles assings a libp2p peerID to each role so they are identifiable
// in the network
type Roles = reason.Roles

// Values maps Parameter.String() to Value
type Values = reason.Values

func compareValues(v1, v2 Values) bool {
	if len(v1) != len(v2) {
		return false
	}
	return reflect.DeepEqual(v1, v2)
}

func compareMessages(m1, m2 Messages) bool {
	if len(m1) != len(m2) {
		return false
	}
	for k, m := range m1 {
		om, found := m2[k]
		if !found {
			return false
		}
		// actions are compared when searching in dict:
		// same action same key and the name check
		if m.Action().Name != m.Action().Name ||
			m.Action().From != m.Action().From ||
			m.Action().To != m.Action().To {
			return false
		}
		if m.InstanceKey() != om.InstanceKey() {
			return false
		}
		if !compareValues(m.Parameters(), om.Parameters()) {
			return false
		}
	}
	return true
}
