package listingitem

import (
	"reflect"
	"strconv"
	"strings"

	i "github.com/14jdelap/thought_machine_take_home/internal"
)

type ListingItem struct {
	Timestamp    int
	UserId       string
	Action       string
	Item         string
	ReservePrice float64
	CloseTime    int
}

// ValidateAndAssign checks if the splitRow has valid inputs for a ListingItem,
// and if valid mutates all of the instance's fields and returns nil.
func (v *ListingItem) ValidateAndAssign(splitRow []string) *i.RowParsingError {
	timestamp, err := strconv.Atoi(splitRow[0])
	if err != nil {
		return &i.RowParsingError{"timestamp", reflect.TypeOf(v)}
	}
	v.Timestamp = timestamp

	_, err = strconv.Atoi(splitRow[1])
	if err != nil {
		return &i.RowParsingError{"userId", reflect.TypeOf(v)}
	}
	v.UserId = splitRow[1]

	action := strings.ToLower(splitRow[2])
	if action != "sell" {
		return &i.RowParsingError{"action", reflect.TypeOf(v)}
	}
	v.Action = action

	item := splitRow[3]
	if len(item) < 1 {
		return &i.RowParsingError{"item", reflect.TypeOf(v)}
	}
	v.Item = item

	reservePrice, err := strconv.ParseFloat(splitRow[4], 64)
	if err != nil {
		return &i.RowParsingError{"reservePrice", reflect.TypeOf(v)}
	}
	v.ReservePrice = reservePrice

	closeTime, err := strconv.Atoi(splitRow[5])
	if err != nil {
		return &i.RowParsingError{"closeTime", reflect.TypeOf(v)}
	}
	v.CloseTime = closeTime

	return nil
}
