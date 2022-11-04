package internal

import (
	"fmt"
	"strings"
)

type auctionResult struct {
	closeTime     int
	item          string
	userId        int
	status        string
	pricePaid     float64
	totalBidCount int
	highestBid    bid
	followUpBid   bid
	lowestBid     bid
	reservePrice  float64
}

type bid struct {
	userId    int
	price     float64
	timestamp int
}

type AuctionHouse struct {
	list               []row
	itemIds            map[string]bool
	itemAuctionResults map[string]*auctionResult
	malformedRows      []string
}

type row interface {
	ValidateAndAssign([]string) error
}

func (a *AuctionHouse) ProcessInputs(inputs []byte) {
	splitInputs := strings.Split(string(inputs), "\n")

	for _, row := range splitInputs {
		splitRow := strings.Split(row, "|")
		rowLength := len(splitRow)

		if rowLength == 1 {
			// Can these 2 rows collapse into 1?
			Heartbeat := &Heartbeat{}
			err := Heartbeat.ValidateAndAssign(splitRow)
			if err == nil {
				a.list = append(a.list, Heartbeat)
			}
		} else if rowLength == 5 {
			BidItem := &BidItem{}
			err := BidItem.ValidateAndAssign(splitRow)
			if err == nil {
				a.list = append(a.list, BidItem)
			}
		} else if rowLength == 6 {
			ListingItem := &ListingItem{}
			err := ListingItem.ValidateAndAssign(splitRow)
			if err == nil {
				a.list = append(a.list, ListingItem)
			}
		}
	}
}

func (a *AuctionHouse) HoldAuction() {
	a.itemAuctionResults = make(map[string]*auctionResult)
	for _, row := range a.list {
		switch row.(type) {
		case *ListingItem:
			a.listItem(row.(*ListingItem))
		case *BidItem:
			a.bidForItem(row.(*BidItem))
		}
	}

	a.reviewSales()
}

func (a *AuctionHouse) listItem(row *ListingItem) {
	_, present := a.itemAuctionResults[row.item]

	if present {
		return
	}

	a.itemAuctionResults[row.item] = &auctionResult{
		closeTime:    row.closeTime,
		item:         row.item,
		reservePrice: row.reservePrice,
	}
}

func (a *AuctionHouse) bidForItem(row *BidItem) {
	currentBid := bid{row.userId, row.bidAmount, row.timestamp}

	auctionItem, present := a.itemAuctionResults[row.item]

	// Do not count bid if the item hasn't been listed
	if !present {
		return
	}

	closeTime := auctionItem.closeTime

	// Do not count bid if the auction for the item has closed
	if closeTime < currentBid.timestamp {
		return
	}

	// If zero value, no bid has been placed
	if auctionItem.lowestBid.price == 0 {
		auctionItem.lowestBid = currentBid
		auctionItem.followUpBid = currentBid
		auctionItem.highestBid = currentBid
	} else {
		if auctionItem.highestBid.price < currentBid.price {
			auctionItem.followUpBid = auctionItem.highestBid
			auctionItem.highestBid = currentBid
		} else if auctionItem.followUpBid.price < currentBid.price {
			auctionItem.followUpBid = currentBid
		} else if auctionItem.lowestBid.price > currentBid.price {
			auctionItem.lowestBid = currentBid
		}
	}
	auctionItem.totalBidCount += 1
}

// reviewSales loops over every auctionResult
func (a *AuctionHouse) reviewSales() {
	for k := range a.itemAuctionResults {
		result := a.itemAuctionResults[k]
		if result.highestBid.price >= result.reservePrice {
			result.status = "SOLD"
			result.pricePaid = result.followUpBid.price
			result.userId = result.highestBid.userId
		} else {
			result.status = "UNSOLD"
		}
	}
}

// Loop over every AuctionResult and log their outcomes
func (a *AuctionHouse) AnnounceResults() {
	fmt.Println("AUCTION RESULTS")
	for k := range a.itemAuctionResults {
		r := a.itemAuctionResults[k]
		message := fmt.Sprintf("%d|%s|%d|%s|%s|%d|%s|%s",
			r.closeTime,
			r.item,
			r.userId,
			r.status,
			fmt.Sprintf("%.2f", r.pricePaid), // Wrong
			r.totalBidCount,
			fmt.Sprintf("%.2f", r.highestBid.price),
			fmt.Sprintf("%.2f", r.lowestBid.price),
		)
		fmt.Println("  ", message)
	}
}
