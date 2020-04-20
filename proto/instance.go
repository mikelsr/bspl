package proto

// Instance of a BSPL protocol
type Instance interface {
	Protocol() Protocol
	Key() string
	Parameters() []InstanceParameter
}

// InstanceParameter associates a value to a Protocol parameter
type InstanceParameter struct {
	Parameter Parameter
	Value string
}
