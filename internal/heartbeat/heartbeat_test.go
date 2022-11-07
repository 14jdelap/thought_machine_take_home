package heartbeat

import (
	"fmt"
	"reflect"
	"testing"

	i "github.com/14jdelap/thought_machine_take_home/internal"
)

func TestValidateAndAssign(t *testing.T) {
	h := &Heartbeat{}
	var tests = []struct {
		description string
		splitRow    []string
		want        *i.RowParsingError
	}{
		{"Happy path: expected output", []string{"12"}, nil},
		{"Happy path: unexpected but valid output", []string{"-37"}, nil},
		{"Unhappy path: string timestamp", []string{"wrong timestamp"}, &i.RowParsingError{"timestamp", reflect.TypeOf(h)}},
		{"Unhappy path: empty timestamp", []string{""}, &i.RowParsingError{"timestamp", reflect.TypeOf(h)}},
	}

	for _, tt := range tests {
		testName := fmt.Sprintf("%v", tt.splitRow)
		t.Run(testName, func(t *testing.T) {
			err := h.ValidateAndAssign(tt.splitRow)
			if err == nil && tt.want == nil {
				return
			} else if err.Column == tt.want.Column && err.RowType == tt.want.RowType {
				return
			} else {
				t.Errorf("got %s, want %s", err, tt.want)
			}
		})
	}
}
