package main

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
)

type ProductListing struct {
	Sku         string
	Description string
	Link        string
	Vendor      string
	InStock     string
	LastInStock string
	Price       string
}

const website string = "https://rpilocator.com/?country=US"

// Gets the table rows
func getRows(ctx context.Context) []*cdp.Node {
	var rows []*cdp.Node
	err := chromedp.Run(ctx,
		chromedp.Nodes(
			"#myTable > tr",
			&rows,
			chromedp.ByQueryAll,
		),
	)

	if err != nil {
		panic(err)
	}

	return rows
}

// Gets the row columns
func getColumns(ctx context.Context, row *cdp.Node) []*cdp.Node {
	var columns []*cdp.Node
	err := chromedp.Run(ctx,
		chromedp.Nodes("td", &columns, chromedp.ByQueryAll, chromedp.FromNode(row)),
	)

	if err != nil {
		panic(err)
	}

	return columns
}

// Save JSON
func saveJson(slice []ProductListing) {
	json, _ := json.MarshalIndent(slice, "", "    ")
	// fmt.Println(string(json))

	outputFile, err := os.Create(`./output/output.json`)
	if err != nil {
		panic(err)
	}

	outputFile.Write(json)
}

// Get all the listings
func getListings(ctx context.Context) {
	var listings []ProductListing
	var rows = getRows(ctx)

	for _, row := range rows {
		var listing = ProductListing{}

		// Get the columns for each row
		var columns = getColumns(ctx, row)

		// Get the info from each column
		err := chromedp.Run(ctx,
			chromedp.Text([]cdp.NodeID{columns[0].NodeID}, &listing.Sku, chromedp.ByNodeID),
			chromedp.Text([]cdp.NodeID{columns[1].NodeID}, &listing.Description, chromedp.ByNodeID),
			chromedp.Text([]cdp.NodeID{columns[4].NodeID}, &listing.Vendor, chromedp.ByNodeID),
			chromedp.Text([]cdp.NodeID{columns[5].NodeID}, &listing.InStock, chromedp.ByNodeID),
			chromedp.Text([]cdp.NodeID{columns[6].NodeID}, &listing.LastInStock, chromedp.ByNodeID),
			chromedp.Text([]cdp.NodeID{columns[7].NodeID}, &listing.Price, chromedp.ByNodeID),
			// Get the purchase link
			chromedp.ActionFunc(func(ctx context.Context) error {
				id, err := dom.QuerySelector(columns[2].NodeID, "a").Do(ctx)
				if err != nil {
					panic(err)
				}

				chromedp.AttributeValue([]cdp.NodeID{id}, "href", &listing.Link, nil, chromedp.ByNodeID).Do(ctx)
				return nil
			}),
		)
		if err != nil {
			panic(err)
		}

		listings = append(listings, listing)
	}

	saveJson(listings)
}

func main() {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Set a 15 second timeout
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// Navigate and wait for the table to populate
	err := chromedp.Run(ctx,
		chromedp.Navigate(website),
		chromedp.WaitNotPresent(`.dataTables_empty`),
		chromedp.ActionFunc(func(ctx context.Context) error {
			getListings(ctx)
			return nil
		}),
	)

	if err != nil {
		panic(err)
	}
}
