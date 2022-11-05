package biditem

import (
	"reflect"
	"strconv"
	"strings"

	i "github.com/14jdelap/thought_machine_take_home/internal"
)

type BidItem struct {
	Timestamp int
	UserId    int
	Action    string
	Item      string
	BidAmount float64
}

func (b *BidItem) ValidateAndAssign(splitRow []string) *i.RowParsingError {
	timestamp, err := strconv.Atoi(splitRow[0])
	if err != nil {
		return &i.RowParsingError{"timestamp", reflect.TypeOf(b)}
	}
	b.Timestamp = timestamp

	userId, err := strconv.Atoi(splitRow[1])
	if err != nil {
		return &i.RowParsingError{"userId", reflect.TypeOf(b)}
	}
	b.UserId = userId

	action := strings.ToLower(splitRow[2])
	if action != "bid" {
		return &i.RowParsingError{"action", reflect.TypeOf(b)}
	}
	b.Action = action

	item := splitRow[3]
	if len(item) < 1 {
		return &i.RowParsingError{"item", reflect.TypeOf(b)}
	}
	b.Item = item

	bidAmount, err := strconv.ParseFloat(splitRow[4], 64)
	if err != nil {
		return &i.RowParsingError{"bidAmount", reflect.TypeOf(b)}
	}
	b.BidAmount = bidAmount

	return nil
}
