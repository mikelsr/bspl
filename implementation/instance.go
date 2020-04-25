package implementation

import (
	"fmt"
	"strings"

	"github.com/mikelsr/bspl"
	"github.com/mikelsr/bspl/proto"
	"github.com/mikelsr/bspl/reason"
)

// Instance of a protocol
type Instance struct {
	protocol proto.Protocol
	roles    Roles
	values   Values
}

// NewInstance is the default constructor for Instance.
func NewInstance(protocol proto.Protocol, roles Roles) *Instance {
	return &Instance{protocol: protocol, roles: roles, values: make(Values)}
}

// Diff identifies what action has been run between two versions of an
// instance. It returns the action, the new values and an error.
// Currently only one action is supported between instace versions.
// An action slice is returned because two actions may have happened,
// e.g. Accept or Reject. In that case the Reasoner must find out which
// one it was.
func (i *Instance) Diff(j reason.Instance) ([]bspl.Action, Values, error) {
	diffValues := make(Values)
	diffParams := make(map[string]bspl.Parameter)
	for paramStr, newValue := range j.Parameters() {
		value, found := i.Parameters()[paramStr]
		if found {
			if value != newValue {
				return nil, nil, fmt.Errorf("Mismatched values for param '%s'", paramStr)
			} else if value != "" {
				continue
			}
		}
		param, found := i.paramFromString(paramStr)
		if !found {
			panic(fmt.Errorf("Parameter not found: '%s'", paramStr))
		}
		diffParams[paramStr] = param
		diffValues[param.Name] = newValue
	}
	// find what actions have happened between the two
	// instances. Currently only one action is supported.
	actions := make([]bspl.Action, 0)
	for _, action := range i.protocol.Actions {
		matches := 0
		for _, out := range action.Outs() {
			for _, param := range diffParams {
				if out.Name == param.Name {
					matches++
				}
			}
		}
		if matches == len(diffParams) {
			actions = append(actions, action)
		}
	}
	if len(actions) == 0 {
		return nil, nil, fmt.Errorf(
			"No action identified for parameters '%s'", diffValues)

	}
	return actions, diffValues, nil
}

// Equals compares two instances.
func (i *Instance) Equals(j reason.Instance) bool {
	if !(i.Key() == j.Key() && i.values.Equals(j.Parameters())) {
		return false
	}
	if len(i.Roles()) != len(j.Roles()) {
		return false
	}
	for k, v1 := range i.Roles() {
		if v2 := i.Roles()[k]; v1 != v2 {
			return false
		}
	}
	return true
}

// GetValue returns the value of the parameter of an instance.
func (i *Instance) GetValue(parameter string) string {
	for _, param := range i.Protocol().Parameters() {
		if param.Name == parameter {
			return i.Parameters()[param.String()]
		}
	}
	return ""
}

// Key of the instance.
func (i *Instance) Key() string {
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

// Parameters of the Instance.
func (i *Instance) Parameters() Values {
	return i.values
}

// Protocol of the Instance.
func (i *Instance) Protocol() proto.Protocol {
	return i.protocol
}

// Roles of the Instance, composed of the protocol name,
// the parameters and the values of the parameters.
func (i *Instance) Roles() Roles {
	return i.roles
}

// SetValue of an instance parameter.
func (i *Instance) SetValue(parameter string, value string) {
	for _, param := range i.Protocol().Parameters() {
		if param.Name == parameter {
			i.Parameters()[param.String()] = value
		}
	}
}

// Update updates an instance given the same instance
// with new actions. WARNING: Currently the instances
// must be updated EACH action, as the search for
// the run action only looks for one.
func (i *Instance) Update(j reason.Instance) error {
	_, values, err := i.Diff(j)
	if err != nil {
		return err
	}
	// Set parameter value
	for k, v := range values {
		i.SetValue(k, v)
	}
	return nil
}

// paramFromString searches for the proto.Parameter struct in an instance
// givenit's string form (in ID key). It would be faster if it where parsed
// but this way we ensure its validity.
func (i *Instance) paramFromString(str string) (proto.Parameter, bool) {
	for _, a := range i.protocol.Actions {
		for _, param := range a.Params {
			if param.String() == str {
				return param, true
			}
		}
	}
	return proto.Parameter{}, false
}
