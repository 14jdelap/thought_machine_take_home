package heartbeat

import (
	"reflect"
	"strconv"

	i "github.com/14jdelap/thought_machine_take_home/internal"
)

type Heartbeat struct {
	Timestamp int
}

func (h *Heartbeat) ValidateAndAssign(splitRow []string) *i.RowParsingError {
	timestamp, err := strconv.Atoi(splitRow[0])
	if err != nil {
		return &i.RowParsingError{"timestamp", reflect.TypeOf(h)}
	}
	h.Timestamp = timestamp

	return nil
}
