package internal

import (
	"errors"
	"strconv"
	"strings"
)

type ListingItem struct {
	timestamp    int
	userId       int
	action       string
	item         string
	reservePrice float64
	closeTime    int
}

func (v *ListingItem) ValidateAndAssign(splitRow []string) error {
	timestamp, err := strconv.Atoi(splitRow[0])
	if err != nil {
		return errors.New("Error when parsing timestamp")
	}
	v.timestamp = timestamp

	userId, err := strconv.Atoi(splitRow[1])
	if err != nil {
		return errors.New("Error when parsing user id")
	}
	v.userId = userId

	action := strings.ToLower(splitRow[2])
	if action != "sell" {
		return errors.New("Error when parsing action")
	}
	v.action = action

	item := splitRow[3]
	if len(item) < 1 {
		return errors.New("Error when parsing item")
	}
	v.item = item

	reservePrice, err := strconv.ParseFloat(splitRow[4], 64)
	if err != nil {
		return errors.New("Error when parsing reserve price")
	}
	v.reservePrice = reservePrice

	closeTime, err := strconv.Atoi(splitRow[5])
	if err != nil {
		return errors.New("Error when parsing close time")
	}
	v.closeTime = closeTime

	return nil
}
