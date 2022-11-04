package internal

import (
	"errors"
	"strconv"
)

type Heartbeat struct {
	timestamp int
}

func (h *Heartbeat) ValidateAndAssign(splitRow []string) error {
	timestamp, err := strconv.Atoi(splitRow[0])
	if err != nil {
		return errors.New("Error when parsing timestamp")
	}
	h.timestamp = timestamp

	return nil
}
