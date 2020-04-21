package instance

import (
	"reflect"
	"strings"

	"github.com/mikelsr/bspl/proto"
)

// Input of a protocol contains a map of the values of the input
// parameters of a protocol and the names of said parameters
// map[name]value
type Input map[string]string

// Roles assings a libp2p peerID to each role so they are identifiable
// in the network
type Roles map[proto.Role]string

// Instance of a protocol
type Instance struct {
	Protocol proto.Protocol
	Roles    Roles
	Values   Values
	Messages map[string]Message
}

// Value of a parameter
type Value string

// Values maps Parameter.String() to Value
type Values map[string]Value

// Message is an instance of an action for an instance of a protocol
type Message struct {
	InstanceKey string
	Action      proto.Action
	Values      Values
}

// Messages maps Instance.Key() to Message
type Messages map[string]Message

// Key of the instance
func (i Instance) Key() string {
	keys := i.Protocol.Keys()
	n := len(keys)

	var sb strings.Builder
	sb.WriteString(i.Protocol.Key())
	sb.WriteRune(instanceSeparator)
	for j, k := range keys {
		v := i.Values[k.String()]
		sb.WriteString(string(v))
		if j != n-1 {
			sb.WriteRune(proto.KeySeparator)
		}
	}
	return sb.String()
}

// Equals compares two instances
func (i Instance) Equals(j Instance) bool {
	// check protocols
	if !reflect.DeepEqual(i.Protocol, j.Protocol) {
		return false
	}
	// check roles
	if !reflect.DeepEqual(i.Roles, j.Roles) {
		return false
	}
	// check values
	if !compareValues(i.Values, j.Values) {
		return false
	}
	return compareMessages(i.Messages, j.Messages)
}

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
		if m.Action.Name != m.Action.Name ||
			m.Action.From != m.Action.From ||
			m.Action.To != m.Action.To {
			return false
		}
		if m.InstanceKey != om.InstanceKey {
			return false
		}
		if !compareValues(m.Values, om.Values) {
			return false
		}
	}
	return true
}
