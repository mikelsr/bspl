package instance

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/mikelsr/bspl/parser"
	"github.com/mikelsr/bspl/proto"
)

type messageMarshaller struct {
	InstanceKey string `json:"instance_key"`
	Action      string `json:"action"`
	Values      Values `json:"values"`
}

type instanceMarshaller struct {
	Protocol string                       `json:"protocol"`
	Roles    Roles                        `json:"roles"`
	Values   Values                       `json:"protocol_values"`
	Messages map[string]messageMarshaller `json:"messages"`
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
		Protocol: i.Protocol.String(),
		Roles:    i.Roles,
		Values:   i.Values,
	}
	im.Messages = make(map[string]messageMarshaller)
	for k, m := range i.Messages {
		mm := messageMarshaller{
			InstanceKey: i.Key(),
			Action:      m.Action.String(),
			Values:      m.Values,
		}
		im.Messages[k] = mm
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
	i.Protocol = p
	i.Roles = im.Roles
	i.Values = im.Values
	i.Messages = make(Messages)
	for k, v := range im.Messages {
		a := new(proto.Action)
		if err = UnmarshalAction(a, []byte(v.Action)); err != nil {
			return err
		}
		m := Message{InstanceKey: v.InstanceKey, Action: *a, Values: v.Values}
		i.Messages[k] = m
	}
	return nil
}

// Marshal a Message
func (m *Message) Marshal() ([]byte, error) {
	mm := messageMarshaller{
		InstanceKey: m.InstanceKey,
		Action:      m.Action.String(),
		Values:      m.Values,
	}
	return json.Marshal(mm)
}

// Unmarshal a Message
func (m *Message) Unmarshal(data []byte) error {
	mm := new(messageMarshaller)
	if err := json.Unmarshal(data, mm); err != nil {
		return err
	}
	m.InstanceKey = mm.InstanceKey
	m.Values = mm.Values
	a := new(proto.Action)
	if err := UnmarshalAction(a, []byte(mm.Action)); err != nil {
		return err
	}
	m.Action = *a
	return nil
}
