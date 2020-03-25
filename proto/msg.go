package proto

// Message representing an action that has taken place in a protocol instance
type Message struct {
	action Action
	From   string
	To     string
	values map[Parameter]string
}

// Keys of a Message
func (m Message) Keys() map[string]string {
	keys := make(map[string]string)
	for _, param := range m.Parameters() {
		if param.Key {
			keys[param.Name] = m.Value(param)
		}
	}
	return keys
}

// Parameters of a message
func (m Message) Parameters() []Parameter {
	return m.action.Parameters()
}

// Value of a Protocol Instance Parameter
func (m Message) Value(p Parameter) string {
	return m.values[p]
}
