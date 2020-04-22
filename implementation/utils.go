package implementation

import (
	"fmt"
	"strings"
)

// emptyProto is used to insert a parameter in a protocol
// string. It is a terrible way of parsing a parameter.
func emptyProto(action ...string) string {
	var sb strings.Builder
	for _, a := range action {
		sb.WriteString(a)
		sb.WriteRune('\n')
	}
	return fmt.Sprintf(`X {
		role A, B
		parameter out C key

		%s
		}`, sb.String())
}
