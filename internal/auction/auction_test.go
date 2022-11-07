package auction

import (
	"fmt"
	"testing"

	b "github.com/14jdelap/thought_machine_take_home/internal/bid_items"
	l "github.com/14jdelap/thought_machine_take_home/internal/listing_item"
)

func TestProcessInputs(t *testing.T) {
	var tests = []struct {
		description string
		inputs      []string
		want        []row
	}{
		{
			"Happy path: all valid",
			[]string{"10|1|SELL|toaster_1|10.00|20", "12|8|BID|toaster_1|7.50", "13|5|BID|toaster_1|12.50"},
			[]row{&l.ListingItem{10, "1", "SELL", "toaster_1", 10.00, 20}, &b.BidItem{12, "8", "BID", "toaster_1", 7.50}, &b.BidItem{13, "5", "BID", "toaster_1", 12.50}},
		},
		{
			"Happy path: some valid",
			[]string{"10|1|SELL|toaster_1|10.00|20", "", "12|8|BID|toaster_1|7.50", "wrong", "13|5|BID|toaster_1|12.50"},
			[]row{&l.ListingItem{10, "1", "SELL", "toaster_1", 10.00, 20}, &b.BidItem{12, "8", "BID", "toaster_1", 7.50}, &b.BidItem{13, "5", "BID", "toaster_1", 12.50}},
		},
		{
			"Unhappy path: none valid",
			[]string{"hi|1|SELL|toaster_1|10.00|20", "12|8|sell|toaster_1|7.50", "13|5|BID||12.50"},
			[]row{},
		},
	}

	for _, tt := range tests {
		testName := fmt.Sprintf("%v", tt.description)
		t.Run(testName, func(t *testing.T) {
			rows := ProcessInputs(tt.inputs)
			if len(rows) != len(tt.want) {
				t.Errorf("got %d rows, want %d rows", len(rows), len(tt.want))
			}
		})
	}
}

func TestHoldAuction(t *testing.T) {
	auctionResult1 := auction{items: make(map[string]*auctionResult)}
	auctionResult1.items["toaster_1"] = &auctionResult{20, "toaster_1", "5", "SOLD", 7.50, 2, bid{"5", 12.50, 13}, bid{"8", 7.50, 12}, bid{"8", 7.50, 12}, 10.00}

	auctionResult2 := auction{items: make(map[string]*auctionResult)}
	auctionResult2.items["toaster_1"] = &auctionResult{20, "toaster_1", "8", "SOLD", 10.00, 2, bid{"8", 10.00, 12}, bid{"8", 10.00, 12}, bid{"8", 10.00, 12}, 10.00}

	auctionResult3 := auction{items: make(map[string]*auctionResult)}
	auctionResult3.items["toaster_1"] = &auctionResult{20, "toaster_1", "8", "SOLD", 10.00, 1, bid{"8", 15.00, 12}, bid{"8", 15.00, 12}, bid{"8", 15.00, 12}, 10.00}

	auctionResult4 := auction{items: make(map[string]*auctionResult)}
	auctionResult4.items["toaster_1"] = &auctionResult{20, "toaster_1", "8", "SOLD", 10.00, 1, bid{"8", 15.00, 12}, bid{"8", 15.00, 12}, bid{"8", 15.00, 12}, 10.00}

	auctionResult5 := auction{items: make(map[string]*auctionResult)}
	auctionResult5.items["toaster_1"] = &auctionResult{20, "toaster_1", "", "UNSOLD", 0.00, 2, bid{"5", 8.50, 13}, bid{"8", 7.50, 12}, bid{"8", 7.50, 12}, 10.00}

	auctionResult6 := auction{items: make(map[string]*auctionResult)}
	auctionResult6.items["toaster_1"] = &auctionResult{20, "toaster_1", "", "UNSOLD", 0.00, 0, bid{}, bid{}, bid{}, 10.00}

	auctionResult7 := auction{items: make(map[string]*auctionResult)}

	auctionResult8 := auction{items: make(map[string]*auctionResult)}
	auctionResult7.items["toaster_1"] = &auctionResult{20, "toaster_1", "", "UNSOLD", 0.00, 0, bid{}, bid{}, bid{}, 10.00}

	var tests = []struct {
		description string
		inputs      []row
		want        auction
	}{
		{
			"Happy path: item is sold under the reserve price because the highest bid was above the reserve",
			[]row{&l.ListingItem{10, "1", "SELL", "toaster_1", 10.00, 20}, &b.BidItem{12, "8", "BID", "toaster_1", 7.50}, &b.BidItem{13, "5", "BID", "toaster_1", 12.50}},
			auctionResult1,
		},
		{
			"Happy path: item with 2 equal bids is sold to the earliest bidder",
			[]row{&l.ListingItem{10, "1", "SELL", "toaster_1", 10.00, 20}, &b.BidItem{12, "8", "BID", "toaster_1", 10.00}, &b.BidItem{13, "5", "BID", "toaster_1", 10.00}},
			auctionResult2,
		},
		{
			"Happy path: item is sold at the reserve price despite higher bid",
			[]row{&l.ListingItem{10, "1", "SELL", "toaster_1", 10.00, 20}, &b.BidItem{12, "8", "BID", "toaster_1", 15.00}},
			auctionResult3,
		},
		{
			"Happy path: auction disregards the user's second, lower bid",
			[]row{&l.ListingItem{10, "1", "SELL", "toaster_1", 10.00, 20}, &b.BidItem{12, "8", "BID", "toaster_1", 15.00}, &b.BidItem{13, "8", "BID", "toaster_1", 9.00}},
			auctionResult4,
		},
		{
			"Unhappy path: item isn't sold because the bid is under the reserve price",
			[]row{&l.ListingItem{10, "1", "SELL", "toaster_1", 10.00, 20}, &b.BidItem{12, "8", "BID", "toaster_1", 7.50}, &b.BidItem{13, "5", "BID", "toaster_1", 8.50}},
			auctionResult5,
		},
		{
			"Unhappy path: no bids",
			[]row{&l.ListingItem{10, "1", "SELL", "toaster_1", 10.00, 20}},
			auctionResult6,
		},
		{
			"Unhappy path: no items on sale",
			[]row{&b.BidItem{12, "8", "BID", "toaster_1", 7.50}, &b.BidItem{13, "5", "BID", "toaster_1", 8.50}},
			auctionResult7,
		},
		{
			"Unhappy path: item not sold because bid arrived after close",
			[]row{&b.BidItem{12, "8", "BID", "toaster_1", 7.50}, &b.BidItem{21, "5", "BID", "toaster_1", 18.50}},
			auctionResult8,
		},
	}

	for _, tt := range tests {
		testName := fmt.Sprintf("%v", tt.description)
		t.Run(testName, func(t *testing.T) {
			auction := HoldAuction(tt.inputs)
			for i := range auction.items {
				result := auction.items[i]

				sameEnd := result.closeTime == tt.want.items[i].closeTime
				if !sameEnd {
					t.Errorf("got %d closeTime, want %d closeTime", result.closeTime, tt.want.items[i].closeTime)
				}

				sameItem := result.item == tt.want.items[i].item
				if !sameItem {
					t.Errorf("got %s item, want %s item", result.item, tt.want.items[i].item)
				}

				sameUserId := result.userId == tt.want.items[i].userId
				if !sameUserId {
					t.Errorf("got %s userId, want %s userId", result.userId, tt.want.items[i].userId)
				}

				sameStatus := result.status == tt.want.items[i].status
				if !sameStatus {
					t.Errorf("got %s status, want %s status", result.status, tt.want.items[i].status)
				}

				samePricePaid := result.pricePaid == tt.want.items[i].pricePaid
				if !samePricePaid {
					t.Errorf("got %.2f pricePaid, want %.2f pricePaid", result.pricePaid, tt.want.items[i].pricePaid)
				}

				sameBidCount := result.totalBidCount == tt.want.items[i].totalBidCount
				if !sameBidCount {
					t.Errorf("got %d totalBidCount, want %d totalBidCount", result.totalBidCount, tt.want.items[i].totalBidCount)
				}

				sameHighestBid := result.highestBid.price == tt.want.items[i].highestBid.price
				if !sameHighestBid {
					t.Errorf("got %.2f highestBid, want %.2f highestBid", result.highestBid.price, tt.want.items[i].highestBid.price)
				}

				sameFollowUpBid := result.followUpBid.price == tt.want.items[i].followUpBid.price
				if !sameFollowUpBid {
					t.Errorf("got %.2f followUpBid, want %.2f followUpBid", result.followUpBid.price, tt.want.items[i].followUpBid.price)
				}

				sameLowestBid := result.lowestBid.price == tt.want.items[i].lowestBid.price
				if !sameLowestBid {
					t.Errorf("got %.2f lowestBid, want %.2f lowestBid", result.lowestBid.price, tt.want.items[i].lowestBid.price)
				}

				sameReservePrice := result.reservePrice == tt.want.items[i].reservePrice
				if !sameReservePrice {
					t.Errorf("got %.2f reservePrice, want %.2f reservePrice", result.reservePrice, tt.want.items[i].reservePrice)
				}
			}
		})
	}
}
