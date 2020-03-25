package proto

// Instance of a Protocol Implementation
type Instance struct {
	protocol Protocol
	messages []Message
	values   map[Parameter]string
}

// Keys of an Instance
func (i Instance) Keys() map[string]string {
	keys := make(map[string]string)
	for _, param := range i.protocol.Parameters() {
		if param.Key {
			keys[param.Name] = i.Value(param)
		}
	}
	return keys
}

// Parameters of an instanced Protocol
func (i Instance) Parameters() []Parameter {
	return i.protocol.Parameters()
}

// Value of a Protocol Instance Parameter
func (i Instance) Value(p Parameter) string {
	return i.values[p]
}
