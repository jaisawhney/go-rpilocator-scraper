package main

import (
	"context"
	"fmt"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

type ScrapedData struct {
	sku         string
	description string
	link        string
	vendor      string
	stockStatus string
	lastInStock string
	price       string
}

const website string = "https://rpilocator.com/?country=US"

func main() {
	var data []ScrapedData

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var rows []*cdp.Node
	err := chromedp.Run(ctx,
		chromedp.Navigate(website),
		chromedp.WaitNotPresent(`.dataTables_empty`), // Wait for the table to populate
		chromedp.Nodes(
			"#myTable > tr",
			&rows,
			chromedp.ByQueryAll,
		),
	)

	if err != nil {
		panic(err)
	}

	for _, row := range rows {
		var columns []*cdp.Node
		err = chromedp.Run(ctx,
			chromedp.Nodes("td", &columns, chromedp.ByQueryAll, chromedp.FromNode(row)),
		)

		if err != nil {
			panic(err)
		}
	}
	fmt.Println(data)
}
