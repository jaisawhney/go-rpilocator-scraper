package main

import (
	"context"

	"github.com/chromedp/chromedp"
)

type ScrapedData struct {
	description string
	link        string
	vendor      string
	stockStatus string
	lastInStock string
	price       string
}

const website string = "https://rpilocator.com/?country=US"

func main() {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.Navigate(website),
		chromedp.WaitNotPresent(`.dataTables_empty`), // Wait for the table to populate
	)

	if err != nil {
		panic(err)
	}
}
