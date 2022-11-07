package biditem

import (
	"fmt"
	"reflect"
	"testing"

	i "github.com/14jdelap/thought_machine_take_home/internal"
)

func TestValidateAndAssign(t *testing.T) {
	b := &BidItem{}
	var tests = []struct {
		description string
		splitRow    []string
		want        *i.RowParsingError
	}{
		// Happy paths
		{"Happy path", []string{"12", "8", "BID", "toaster_1", "7.50"}, nil},
		// Unhappy paths
		{"Unhappy path: string timestamp", []string{"wrong timestamp", "8", "BID", "toaster_1", "7.50"}, &i.RowParsingError{"timestamp", reflect.TypeOf(b)}},
		{"Unhappy path: string userId", []string{"12", "wrong id", "BID", "toaster_1", "7.50"}, &i.RowParsingError{"userId", reflect.TypeOf(b)}},
		{"Unhappy path: action isn't bid", []string{"12", "8", "SELL", "toaster_1", "7.50"}, &i.RowParsingError{"action", reflect.TypeOf(b)}},
		{"Unhappy path: empty item", []string{"12", "8", "BID", "", "7.50"}, &i.RowParsingError{"item", reflect.TypeOf(b)}},
		{"Unhappy path: empty bidAmount", []string{"12", "8", "BID", "T", ""}, &i.RowParsingError{"bidAmount", reflect.TypeOf(b)}},
		{"Unhappy path: string bidAmount", []string{"12", "8", "BID", "T", "hello"}, &i.RowParsingError{"bidAmount", reflect.TypeOf(b)}},
	}

	for _, tt := range tests {
		testName := fmt.Sprintf("%v", tt.splitRow)
		t.Run(testName, func(t *testing.T) {
			err := b.ValidateAndAssign(tt.splitRow)
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
