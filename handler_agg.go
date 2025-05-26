package main

import (
	"context"
	"fmt"
	"log"
	"time"
)


func handlerAgg(s *state, cmd command) error {

	// Check for absence of args
	if len(cmd.Arguments) != 1 {
		return fmt.Errorf("usage: %s <time_between_reqs>", cmd.Name)
	}

	// Parse the time arg
	reqsTime, errParseTime := time.ParseDuration(cmd.Arguments[0])
	if errParseTime != nil {
		return fmt.Errorf("couldn't parse time_between_reqs argument: %w", errParseTime)
	}

	// Run scrapping according to time_between_reqs arg
	ticker := time.NewTicker(reqsTime)
	fmt.Printf("Collecting feeds every %s!\n\n", cmd.Arguments[0])
	for ; ; <-ticker.C {
		fmt.Println("INITIATING FEED")
		fmt.Println("==========================================================")
		scrapeFeeds(s)
		fmt.Println("ENDING FEED")
		fmt.Println("==========================================================")
		fmt.Printf("\n\n")
	}
}

func scrapeFeeds(s *state) {

	dbQueries := s.db
	// Get next feed to fetch based on date
	nextFeed, errGetNextFeed := dbQueries.GetNextFeedToFetch(context.Background())
	if errGetNextFeed != nil {
		log.Println("couldn't get next feed", errGetNextFeed)
		return 
	}

	// Mark fetched feed
	errMarkFeed := dbQueries.MarkFeedFetched(context.Background(), nextFeed.ID)
	if errMarkFeed != nil {
		log.Println("couldn't mark fetched feed", errMarkFeed)
		return
	}

	// Get feed by URL
	feed, errGetFeedByURL := fetchFeed(context.Background(), nextFeed.Url)
	if errGetFeedByURL != nil {
		log.Println("couldn't fetch feed by URL", errGetFeedByURL)
		return 
	}

	printRSSFeed(feed)
}
func printRSSFeed(feed *RSSFeed) {

	// Print channel's info
	fmt.Printf(" * Title: %v\n", feed.Channel.Title)
	for i, link := range feed.Channel.Link {
		if link.Text != "" {
			fmt.Printf(" * Link %v: %v\n", i+1, link.Text)
		} else if link.Href != "" {
			fmt.Printf(" * Link %v: %v\n", i+1, link.Href)
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
