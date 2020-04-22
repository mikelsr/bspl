package implementation

import "github.com/mikelsr/bspl/proto"

// Message is an instance of an action for an instance of a protocol
type Message struct {
	instanceKey string
	action      proto.Action
	values      Values
}

// NewMessage is the default constructor of Message
func NewMessage(instanceKey string, action proto.Action, values Values) Message {
	return Message{instanceKey: instanceKey, action: action, values: values}
}

// InstanceKey of the Message
func (m Message) InstanceKey() string {
	return m.instanceKey
}

// Action of the Message
func (m Message) Action() proto.Action {
	return m.action
}

// Parameters of the Message
func (m Message) Parameters() Values {
	return m.values
}
