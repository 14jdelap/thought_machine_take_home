# Auction House: Thought Machine Take Home

From the **top-level directory** you can:

- Run this program with `go run ./... [path=string]`
  - If no path is given, the default of `inputs/input.txt` is used
- Run the tests with `go test ./...`

## Testing edge cases

Besides running the tests, `inputs/malformed_inputs.txt` shows how the program handles edge cases. It has the same result as `inputs/input.txt`, but has many rows that are incorrectly formatted.

## Assumptions

- The file will have data in ascending order by timestamp value (per the instructions)
- Inputs may be syntactically malformed (e.g. empty data, wrong types, empty row)
- If there are two bids and the highest is above the reserve price but the lowest is lower, then the buyer purchases the item with the lower price (i.e. under the reserve price)
  - Why: the instructions only say that at least one bid has to be larger than the reserve price

## Design decisions

### Packages and code organization

I thought of this program throughout as being about 4 entities:

- An auction or auction house that acted as an orchestrator of the auction
- Listing items that represented a user's actions to list items
- Bids for items when a user bid for an existing item
- Heartbeats to increase the time

This is why, `main.go` aside, I separated the code into packages distinctly split up between those 4.

### OOP, state management, and decoupling

Coming from an OOP background, my first iteration leaned much more towards an OOP approach. The biggest difference is in `auction.go`'s data and functions, which were initially much more coupled.

To conform to Go's spirit of high decoupling, I changed almost all methods to functions. The exception are the `ValidateAndAssign` methods. The reason I kept them is that by being methods I can both check the validity of the string data and assign those values to the new struct, instead of having to iterate twice (once to check and return a `bool` and a second time to create the struct with those values).

### Error handling and wrong data

Throughout the program I err on the side of using errors to inform what data should and should not be used to construct the final results of the auction.

For example, the three `ValidateAndAssign` methods return an error when their inputs don't conform to what's expected from their structs. What `ProcessInputs` does with those an error is simply skip that iteration and continue to the next element in the `for` loop.

The only place where I terminate the program is in `main.go` when reading the file's contents return an error. In that situation, an error is most appropriate because there's no auction to hold because of a lack of data.

That's also why the program runs until the end when an input has no valid rows or is an empty text file. Instead of stopping the program, I inform the user simply that no item was listed for sale in the auction. Similarly, this is why I do not use `Heartbeat` instances besides validating that their data is currently formed (i.e., that the timestamp is a valid int) — I can ignore them and I still get the expected results.

## Areas of improvement

### Testing

The biggest improvement to testing would be a better way to handle data.

I'm using many structs and slices, which makes the data I use for testing hard to read (see e.g. `TestHoldAuction` in `auction_test.go`). A big improvement would be writing the data in a separate file and injecting it in its appropriate location for testing. I did not have the time to research/prototype how that could look like, but I expect there must be a better way to handle this.

A more minor topic is testing coverage because it isn't 100%. This isn't a big problem because all the highest impact functions are covered (e.g. `HoldAuction`), but having 100% coverage would make it easier to catch bugs.

### Handling no valid rows

I currently handle this, like mentioned above, by logging to the console in `AnnounceResults`that `"No item was listed for sale"`.

However, this could improve by informing the user if i) the input file was empty or ii) had data for it didn't lead to an item listing (suggesting the file has an input problem).

### Custom errors

I currently use a fairly generic `RowParsingError` in the `ValidateAndAssign` methods. I attempted to prototype more specific errors for each package (e.g. `BidError`, `HeartbeatError`) but wasn't able to finish implementing them because of en error raised that I wasn't able to solve.

However, this is a small improvement that I would adopt if I were to continue working on this.
