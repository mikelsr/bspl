package reason

import "github.com/mikelsr/bspl/proto"

// Message instances an action
type Message interface {
	// Action a message belongs to
	Action() proto.Action
	// InstanceKey returns the key of the Instance
	// the Message belongs to
	InstanceKey() string
	// Marshal a Message to bytes
	Marshal() ([]byte, error)
	// Parameters of the Message
	Parameters() Values
	// Unmarshal a Message from bytes
	Unmarshal([]byte) (Message, error)
}

// Instance of a Protocol
type Instance interface {
	// AddMessage should check if a requirements for adding
	// a message to the instance are met and if so add it
	AddMessage(Message) error
	// Equals should compare two instances
	Equals(Instance) bool
	// Key of the Instance
	Key() string
	// Marshal an Instance to bytes
	Marshal() ([]byte, error)
	// Messages of the Instance
	Messages() Messages
	// Roles of the Instance
	Roles() Roles
	// Parameters of the Instance
	Parameters() Values
	// Protocol of the Instance
	Protocol() proto.Protocol
	// Unmarshal an Instance from bytes
	Unmarshal([]byte) (Instance, error)
}
