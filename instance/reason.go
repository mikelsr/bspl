package instance

import "github.com/mikelsr/bspl/proto"

// Reasoner handles the protocol instances and actions related to them
type Reasoner interface {
	// Instantiate a protocol. Check if the assigned role is a role
	// the reasoner is willing to play.
	Instantiate(p proto.Protocol, ins Input) (Instance, error)
	// NewMessage creates a message for an action of an instance
	// for that it must check the validity of the message in the
	// current state of the instance.
	NewMessage(i Instance, a proto.Action) (Message, error)
	// Addmesage adds a message to the instance if the message was
	// created by another role.
	Addmesage(i Instance, m Message) (Instance, error)
}
