package main

import (
	"context"
	"fmt"
)

const (
	FETCH_URL = "https://www.wagslane.dev/index.xml"
)

func handlerAgg(s *state, cmd command) error {

	// Check for absence of args
	if len(cmd.Arguments) != 0 {
		return fmt.Errorf("%s doesn't accept arguments", cmd.Name)
	}

	// Fetch feed
	rssFeed, errFetch := fetchFeed(context.Background(), FETCH_URL)
	if errFetch != nil {
		return fmt.Errorf("couldn't fetch given URL: %w", errFetch)
	}

	// Print feed
	printRSSFeed(rssFeed)
	return nil
}

func printRSSFeed(feed *RSSFeed) {

	// Print channel's info
	fmt.Printf(" * Title: %v\n", feed.Channel.Title)
	for i, link := range feed.Channel.Link {
		if link.Text != "" {
			fmt.Printf(" * Link %v: %v\n", i + 1, link.Text)
		} else if link.Href != "" {
			fmt.Printf(" * Link %v: %v\n", i + 1, link.Href)
		}
	}
	fmt.Printf(" * Description: %v\n", feed.Channel.Description)

	// Print item's list
	fmt.Println(" * Item's list:")
	fmt.Println("--------------------")
	for i, item := range feed.Channel.Item {
		fmt.Printf("* Item %v:\n\n", i+1)
		printRSSItem(item)
		fmt.Println("--------------------")
	}
}

func printRSSItem(item RSSItem) {

	// Print RSSItem's fields
	fmt.Printf("Title: %v\n\n", item.Title)
	fmt.Printf("Link: %v\n\n", item.Link)
	fmt.Printf("Description: %v\n\n", item.Description)
	fmt.Printf("PubDate: %v\n", item.PubDate)
}
