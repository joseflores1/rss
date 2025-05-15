package main

import (
	"context"
	"fmt"
)

const (
	FETCH_URL = "https://www.wagslane.dev/index.xml"
)

func handlerAgg(s *state, cmd command) error {

	if len(cmd.Arguments) != 0 {
		return fmt.Errorf("%s doesn't accept arguments", cmd.Name)
	}

	rssFeed, errFetch := fetchFeed(context.Background(), FETCH_URL)
	if errFetch != nil {
		return fmt.Errorf("couldn't fetch given URL: %w", errFetch)
	}

	printRSSFeed(rssFeed)
	return nil
}

func printRSSFeed(feed *RSSFeed) {

	fmt.Printf(" * Title: %v\n", feed.Channel.Title)
	fmt.Printf(" * Link: %v\n", feed.Channel.Link)
	fmt.Printf(" * Description: %v\n", feed.Channel.Description)

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
