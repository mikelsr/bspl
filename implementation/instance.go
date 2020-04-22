package implementation

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/mikelsr/bspl/proto"
)

// Instance of a protocol
type Instance struct {
	protocol proto.Protocol
	roles    Roles
	values   Values
	messages Messages
}

// NewInstance is the default constructor for Instance
func NewInstance(protocol proto.Protocol, roles Roles, parameters Values) Instance {
	return Instance{protocol: protocol, roles: roles, values: parameters}
}

// Protocol of the Instance
func (i Instance) Protocol() proto.Protocol {
	return i.protocol
}

// Roles of the Instance
func (i Instance) Roles() Roles {
	return i.roles
}

// Parameters of the Instance
func (i Instance) Parameters() Values {
	return i.values
}

// Messages of the Instance
func (i Instance) Messages() Messages {
	return i.messages
}

// AddMessage adds a new message to an instance
func (i Instance) AddMessage(m Message) error {
	if _, found := i.messages[m.Action().String()]; found {
		return fmt.Errorf("Message already exists for action '%s'", m.Action())
	}
	i.messages[m.Action().String()] = m
	return nil
}

// Key of the instance
func (i Instance) Key() string {
	keys := i.protocol.Keys()
	n := len(keys)

	var sb strings.Builder
	sb.WriteString(i.protocol.Key())
	sb.WriteRune(instanceSeparator)
	for j, k := range keys {
		v := i.values[k.String()]
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
	if !reflect.DeepEqual(i.Protocol(), j.Protocol()) {
		return false
	}
	// check roles
	if !reflect.DeepEqual(i.Roles(), j.Roles()) {
		return false
	}
	// check values
	if !i.Parameters().Equals(j.Parameters()) {
		return false
	}
	return i.messages.Equals(j.messages)
}
