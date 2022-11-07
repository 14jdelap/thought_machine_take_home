package heartbeat

import (
	"reflect"
	"strconv"

	i "github.com/14jdelap/thought_machine_take_home/internal"
)

type Heartbeat struct {
	Timestamp int
}

// ValidateAndAssign checks if the splitRow has valid inputs for a Heartbeat,
// and if valid mutates all of the instance's fields and returns nil.
func (h *Heartbeat) ValidateAndAssign(splitRow []string) *i.RowParsingError {
	timestamp, err := strconv.Atoi(splitRow[0])
	if err != nil {
		return &i.RowParsingError{"timestamp", reflect.TypeOf(h)}
	}
	h.Timestamp = timestamp

	return nil
}
