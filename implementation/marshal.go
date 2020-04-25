package implementation

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/mikelsr/bspl/parser"
	"github.com/mikelsr/bspl/proto"
)

type instanceMarshaller struct {
	Protocol string `json:"protocol"`
	Roles    Roles  `json:"roles"`
	Values   Values `json:"protocol_values"`
}

// MarshalAction marshals an Action into bytes
func MarshalAction(a proto.Action) []byte {
	return []byte(a.String())
}

// UnmarshalAction unmarshals an Action from bytes
func UnmarshalAction(a *proto.Action, b []byte) error {
	// To reuse Parse function actions are wrapped in a
	// meaningless protocol, the protocol is parsed.
	wrapper := bytes.NewReader([]byte(emptyProto(string(b))))
	p, err := parser.Parse(wrapper)
	// In correct cases error will be either nil or a ValidationError
	if err != nil {
		switch err.(type) {
		case proto.ValidationError:
			break
		default:
			return err
		}
	}
	if len(p.Actions) == 0 {
		return errors.New("Error unarshalling action")
	}
	*a = p.Actions[0]
	return nil
}

// Marshal an Instance
func (i *Instance) Marshal() ([]byte, error) {
	im := instanceMarshaller{
		Protocol: i.protocol.String(),
		Roles:    i.roles,
		Values:   i.values,
	}
	return json.Marshal(im)
}

// Unmarshal an instance
func (i *Instance) Unmarshal(data []byte) error {
	im := new(instanceMarshaller)
	if err := json.Unmarshal(data, im); err != nil {
		return err
	}
	p, err := parser.Parse(bytes.NewReader([]byte(im.Protocol)))
	if err != nil {
		return err
	}
	i.protocol = p
	i.roles = im.Roles
	i.values = im.Values
	return nil
}
