package listingitem

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/14jdelap/thought_machine_take_home/internal"
)

func TestValidateAndAssign(t *testing.T) {
	l := &ListingItem{}
	var tests = []struct {
		description string
		splitRow    []string
		want        *internal.RowParsingError
	}{
		{"Happy path: expected inputs", []string{"12", "8", "SELL", "toaster_1", "7.50", "20"}, nil},
		{"Happy path: negative timestamp", []string{"-12", "8", "sell", "toaster_1", "7.50", "20"}, nil},
		{"Happy path: irregular action casing", []string{"12", "8", "SeLl", "toaster_1", "750", "20"}, nil},
		{"Unhappy path: string timestamp", []string{"12T", "8", "SeLl", "toaster_1", "750", "20"}, &internal.RowParsingError{"timestamp", reflect.TypeOf(l)}},
		{"Unhappy path: empty timestamp", []string{"", "8", "SeLl", "toaster_1", "750", "20"}, &internal.RowParsingError{"timestamp", reflect.TypeOf(l)}},
		{"Unhappy path: empty userId", []string{"12", "", "SeLl", "toaster_1", "750", "20"}, &internal.RowParsingError{"userId", reflect.TypeOf(l)}},
		{"Unhappy path: string userId", []string{"12", "T", "SeLl", "toaster_1", "750", "20"}, &internal.RowParsingError{"userId", reflect.TypeOf(l)}},
		{"Unhappy path: action isn't sell", []string{"12", "8", "BID", "toaster_1", "750", "20"}, &internal.RowParsingError{"action", reflect.TypeOf(l)}},
		{"Unhappy path: empty action", []string{"12", "8", "", "toaster_1", "750", "20"}, &internal.RowParsingError{"action", reflect.TypeOf(l)}},
		{"Unhappy path: empty item", []string{"12", "8", "SeLl", "", "750", "20"}, &internal.RowParsingError{"item", reflect.TypeOf(l)}},
		{"Unhappy path: string reservePrice", []string{"12", "8", "SeLl", "toaster_1", "7.0.s", "20"}, &internal.RowParsingError{"reservePrice", reflect.TypeOf(l)}},
		{"Unhappy path: empty closeTime", []string{"12", "8", "SeLl", "toaster_1", "750", ""}, &internal.RowParsingError{"closeTime", reflect.TypeOf(l)}},
	}

	for _, tt := range tests {
		testName := fmt.Sprintf("%v", tt.description)
		t.Run(testName, func(t *testing.T) {
			err := l.ValidateAndAssign(tt.splitRow)
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
