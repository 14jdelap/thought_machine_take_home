# thought_machine_take_home

## Problem

- Create an auction house that sells items based on bids
- Conditions for sale
  - There must be one or more bids meeting or in exceess of the reserve price
    - If so, the highest bidder wins the item but pays the price of the 2nd highest bid
      - If two bids are received for the same amount, they pay that price but the earliest bid wins
    - Exception: if there's only one bid they'll pay the reserve price (assuming their bid was equal or higher)
- Conditions for valid bid
  - Arrives after the auction start time and before or ON the closing time
  - Is larger than any previous valid bids submitted by the user
    - Assumption: the *same* user, so a bid might be the highest from a user but not the highest overall if it came from other users
- Heartbeat messages: increase the time periodically to close the auction if bids don't come in
  - Why:in case few bids we need a mechanism to increase time and close an item
- Input: file with auction instructions
  - Need to determine the type of input based on its shape
    - Listing item for sale: `timestamp|user_id|action|item|reserve_price|close_time`
    - Bids on item: `timestamp|user_id|action|item|bid_amount`
    - Heartbeat messages: `timestamp`
- Output at the end of the auction: each item's
  - `close_time|item|user_id|status|price_paid|total_bid_count|highest_bid|lowest_bid`

## Assumptions

- Inputs may be malformed
  - Row might not conform to the shapes described above
  - Row may contain data of different types (e.g. `timestamp` expects `int` that increases monotonically but could receive a lower `int[] or a `string`)
  - Row may be empty! Skip if so.
- Cannot use libraries outside the Go standard library
- User supplied path to the input file: does it have to be an absolute or relative path???
- Input cannot start with a heartbeat
- Only one type of currency
- Order of row outputs is irrelevant
- Timestamps are ordered correctly

## Requirements

- Tests
- Comments
- Error handling of edge cases

## Steps

- Read the text file from the command line
- Process the inputs to create a slice of string slices
  - String slices should not have a type, because they could be either heartbeats, bids on items, or users listing items for sale
  - Enforce rules
    - Maybe have the slice be in a data structure that keeps track of the most recent timestamp to avoid future entries having a lower or repeated value?
    - Have a map that stores item unique ids to avoid repeat item ids?
      - Maybe output message when there's a repeat number id to inform the user of a wrong input?
      - Only add items if it comes from a row with `SELL` action
        - E.g. avoid situations where the row has an item that hasn't appeared before in a `SELL`
    - Maybe store/output all malformed rows?
  - Have a 2nd map of items that go to auction with their result, i.e.:
    - `close_time|item|user_id|status|price_paid|total_bid_count|highest_bid|lowest_bid`
      - Great candidate for a struct type
- After creating the slice only with valid rows, process the rows to create the outputs
  - Go over each row
    - Do not consider bids of non-existent or 


## To do

- Map edge cases with new lists of text files
- Create unit and end-to-end tests (latter might be easier?)
- Consider how concurrency could help, and how to implement it