package reason

import "github.com/mikelsr/bspl/proto"

// Instance of a Protocol
type Instance interface {
	// Diff identifies what action has been run between two versions of an
	// instance. It returns the action, the new values and an error.
	// Currently only one action is supported between instace versions.
	// An action slice is returned because two actions may have happened,
	// e.g. Accept or Reject. In that case the Reasoner must find out which
	// one it was.
	Diff(Instance) ([]proto.Action, Values, error)
	// Equals compares two instances.
	Equals(Instance) bool
	// GetValue returns the value of the parameter of an instance.
	GetValue(string) string
	// Key of the Instance.
	Key() string
	// Marshal an Instance to bytes.
	Marshal() ([]byte, error)
	// Parameters of the Instance.
	Parameters() Values
	// Protocol of the Instance.
	Protocol() proto.Protocol
	// Roles of the Instance.
	Roles() Roles
	// SetValue of an instance parameter.s
	SetValue(string, string)
	// Unmarshal an Instance from bytes
	Unmarshal([]byte) error
	// Update updates an instance given the same instance
	// with new actions. WARNING: Currently the instances
	// must be updated EACH action, as the search for
	// the run action only looks for one.
	Update(Instance) error
}

// Reasoner handles the protocol instances and actions related to them
type Reasoner interface {
	// DropInstance cancels an Instance for whatever motive
	DropInstance(instanceKey string, motive string) error
	// GetInstance returns an Instance given the instance key
	GetInstance(instanceKey string) (Instance, bool)
	// All instances of a Protocol
	Instances(p proto.Protocol) []Instance
	// Instantiate a protocol. Check if the assigned role is a role
	// the reasoner is willing to play.
	Instantiate(p proto.Protocol, roles Roles, ins Values) (Instance, error)
	// RegisterInstance registers an Instance created by another Reasoner
	RegisterInstance(i Instance) error
	// UpdateInstance updates an instance with a newer version of itself
	// as long as a valid run from one to the other.
	UpdateInstance(newVersion Instance) error
}
