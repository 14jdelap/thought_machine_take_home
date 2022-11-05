package auctionhouse

import (
	"fmt"
	"strings"

	i "github.com/14jdelap/thought_machine_take_home/internal"
	b "github.com/14jdelap/thought_machine_take_home/internal/bid_items"
	h "github.com/14jdelap/thought_machine_take_home/internal/heartbeat"
	l "github.com/14jdelap/thought_machine_take_home/internal/listing_item"
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
	itemAuctionResults map[string]*auctionResult
}

type row interface {
	ValidateAndAssign([]string) *i.RowParsingError
}

func (a *AuctionHouse) ProcessInputs(inputs []string) {
	for _, row := range inputs {
		splitRow := strings.Split(row, "|")
		rowLength := len(splitRow)

		if rowLength == 1 {
			// Can these 2 rows collapse into 1?
			heartbeat := &h.Heartbeat{}
			err := heartbeat.ValidateAndAssign(splitRow)
			if err == nil {
				a.list = append(a.list, heartbeat)
			}
		} else if rowLength == 5 {
			bidItem := &b.BidItem{}
			err := bidItem.ValidateAndAssign(splitRow)
			if err == nil {
				a.list = append(a.list, bidItem)
			}
		} else if rowLength == 6 {
			listingItem := &l.ListingItem{}
			err := listingItem.ValidateAndAssign(splitRow)
			if err == nil {
				a.list = append(a.list, listingItem)
			}
		}
	}
}

func (a *AuctionHouse) HoldAuction() {
	a.itemAuctionResults = make(map[string]*auctionResult)
	for _, row := range a.list {
		switch row.(type) {
		case *l.ListingItem:
			a.listItem(row.(*l.ListingItem))
		case *b.BidItem:
			a.bidForItem(row.(*b.BidItem))
		}
	}

	a.reviewSales()
}

func (a *AuctionHouse) listItem(row *l.ListingItem) {
	_, present := a.itemAuctionResults[row.Item]

	if present {
		return
	}

	a.itemAuctionResults[row.Item] = &auctionResult{
		closeTime:    row.CloseTime,
		item:         row.Item,
		reservePrice: row.ReservePrice,
	}
}

func (a *AuctionHouse) bidForItem(row *b.BidItem) {
	currentBid := bid{row.UserId, row.BidAmount, row.Timestamp}

	auctionItem, present := a.itemAuctionResults[row.Item]

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
