package internal

import (
	"errors"
	"strconv"
	"strings"
)

type BidItem struct {
	timestamp int
	userId    int
	action    string
	item      string
	bidAmount float64
}

func (b *BidItem) ValidateAndAssign(splitRow []string) error {
	timestamp, err := strconv.Atoi(splitRow[0])
	if err != nil {
		return errors.New("Error when parsing timestamp")
	}
	b.timestamp = timestamp

	userId, err := strconv.Atoi(splitRow[1])
	if err != nil {
		return errors.New("Error when parsing user id")
	}
	b.userId = userId

	action := strings.ToLower(splitRow[2])
	if action != "bid" {
		return errors.New("Error when parsing action")
	}
	b.action = action

	item := splitRow[3]
	if len(item) < 1 {
		return errors.New("Error when parsing item")
	}
	b.item = item

	bidAmount, err := strconv.ParseFloat(splitRow[4], 64)
	if err != nil {
		return errors.New("Error when parsing bid amount")
	}
	b.bidAmount = bidAmount

	return nil
}
