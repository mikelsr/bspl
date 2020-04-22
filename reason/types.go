package reason

import (
	"reflect"

	"github.com/mikelsr/bspl/proto"
)

// Messages maps action identifiers (.String()) to messages
type Messages map[string]Message

// Values maps component (Protocol, Action...) identifiers to
// string values
type Values map[string]string

// Roles maps protocol roles to string IDs of the agents
// assuming the roles
type Roles map[proto.Role]string

// Equals compares two Messages (Messages = []Message)
func (m Messages) Equals(others Messages) bool {
	if len(m) != len(others) {
		return false
	}
	for k, message := range m {
		om, found := others[k]
		if !found {
			return false
		}
		// actions are compared when searching in dict:
		// same action same key and the name check
		if message.Action().Name != om.Action().Name ||
			message.Action().From != om.Action().From ||
			message.Action().To != om.Action().To {
			return false
		}
		if message.InstanceKey() != om.InstanceKey() {
			return false
		}
		if message.Parameters().Equals(om.Parameters()) {
			return false
		}
	}
	return true
}

// Equals compares two Values (Values = map[string]string)
func (v Values) Equals(others Values) bool {
	if len(v) != len(others) {
		return false
	}
	return reflect.DeepEqual(v, others)
}
