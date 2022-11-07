package auction

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
	userId        string
	status        string
	pricePaid     float64
	totalBidCount int
	highestBid    bid
	followUpBid   bid
	lowestBid     bid
	reservePrice  float64
}

type bid struct {
	userId    string
	price     float64
	timestamp int
}

type auction struct {
	items          map[string]*auctionResult
	userHighestBid map[string]float64
}

type row interface {
	ValidateAndAssign([]string) *i.RowParsingError
}

// ProcessInputs uses the split inputs to create a slice with
// syntactically correct (though not necessarily valid) rows.
func ProcessInputs(inputs []string) []row {
	rows := []row{}
	for _, row := range inputs {
		splitRow := strings.Split(row, "|")
		rowLength := len(splitRow)

		if rowLength == 1 {
			heartbeat := &h.Heartbeat{}
			err := heartbeat.ValidateAndAssign(splitRow)
			if err == nil {
				rows = append(rows, heartbeat)
			}
		} else if rowLength == 5 {
			bidItem := &b.BidItem{}
			err := bidItem.ValidateAndAssign(splitRow)
			if err == nil {
				rows = append(rows, bidItem)
			}
		} else if rowLength == 6 {
			listingItem := &l.ListingItem{}
			err := listingItem.ValidateAndAssign(splitRow)
			if err == nil {
				rows = append(rows, listingItem)
			}
		}
	}
	return rows
}

// HoldAuction iterates over every row and evaluates if a row
// will be used in the auction through listItem and bidForItem.
func HoldAuction(rows []row) auction {
	a := auction{}
	a.items = make(map[string]*auctionResult)
	a.userHighestBid = make(map[string]float64)
	for _, row := range rows {
		switch row.(type) {
		case *l.ListingItem:
			a = listItem(a, row.(*l.ListingItem))
		case *b.BidItem:
			a = bidForItem(a, row.(*b.BidItem))
		}
	}

	return reviewSales(a)
}

// listItem determines a ListingItem's validity, and if so
// uses it in the auction by adding its auctionResult.
func listItem(a auction, row *l.ListingItem) auction {
	_, present := a.items[row.Item]

	// If the same item exits, do not list it twice
	if present {
		return a
	}

	a.items[row.Item] = &auctionResult{
		closeTime:    row.CloseTime,
		item:         row.Item,
		reservePrice: row.ReservePrice,
	}
	return a
}

// bidForItem determines a BidItem's validity, and if valid uses it in the
// auction by mutating the state of the auctionResult.
func bidForItem(a auction, row *b.BidItem) auction {
	currentBid := bid{row.UserId, row.BidAmount, row.Timestamp}

	auctionItem, present := a.items[row.Item]

	// Dismiss bid if the item hasn't been listed
	if !present {
		return a
	}

	closeTime := auctionItem.closeTime

	// Dismiss bid if the auction for the item has closed
	if closeTime < currentBid.timestamp {
		return a
	}

	// Dismiss the bid if the user has a previous, smaller bid
	lastBid, present := a.userHighestBid[currentBid.userId]
	if lastBid > currentBid.price {
		return a
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

	a.userHighestBid[currentBid.userId] = currentBid.price
	auctionItem.totalBidCount += 1

	return a
}

// reviewSales determines what was sold by looping over every AuctionResult
// and updating the state of the auctionResult instances.
func reviewSales(a auction) auction {
	for k := range a.items {
		result := a.items[k]
		if result.highestBid.price >= result.reservePrice {
			result.status = "SOLD"
			result.userId = result.highestBid.userId
			if result.totalBidCount > 1 {
				result.pricePaid = result.followUpBid.price
			} else {
				result.pricePaid = result.reservePrice
			}
		} else {
			result.status = "UNSOLD"
		}
	}
	return a
}

// AnnounceResults loops over every auctionResult and logs their data.
func AnnounceResults(a auction) {
	fmt.Println("AUCTION RESULTS")

	if len(a.items) == 0 {
		fmt.Println("  No item was listed for sale")
		return
	}

	for k := range a.items {
		r := a.items[k]
		message := fmt.Sprintf("%d|%s|%s|%s|%s|%d|%s|%s",
			r.closeTime,
			r.item,
			r.userId,
			r.status,
			fmt.Sprintf("%.2f", r.pricePaid),
			r.totalBidCount,
			fmt.Sprintf("%.2f", r.highestBid.price),
			fmt.Sprintf("%.2f", r.lowestBid.price),
		)
		fmt.Println(" ", message)
	}
}
