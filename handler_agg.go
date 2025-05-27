package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/joseflores1/rss/internal/database"
)

func handlerAgg(s *state, cmd command) error {

	// Check for valid command input
	if len(cmd.Arguments) != 1 {
		return fmt.Errorf("usage: %s <time_between_reqs>", cmd.Name)
	}

	// Parse the time argument
	reqsTime, errParseTime := time.ParseDuration(cmd.Arguments[0])
	if errParseTime != nil {
		return fmt.Errorf("couldn't parse <time_between_reqs> argument: %w", errParseTime)
	}

	// Run scrapping according to time_between_reqs argument
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

func printPost(post database.Post) {

	// Print RSSItem's fields
	fmt.Printf("Title: %v\n\n", post.Title)
	fmt.Printf("Link: %v\n\n", post.Url)
	fmt.Printf("Description: %v\n\n", post.Description)
	fmt.Printf("PubDate: %v\n", post.PublishedAt)
}

func scrapeFeeds(s *state) {

	// Define variables
	dbQueries := s.db

	// Get next feed to fetch based on date
	nextFeed, errGetNextFeed := dbQueries.GetNextFeedToFetch(context.Background())
	if errGetNextFeed != nil {
		log.Println("couldn't get next feed:", errGetNextFeed)
		return
	}

	// Mark fetched feed
	errMarkFeed := dbQueries.MarkFeedFetched(context.Background(), nextFeed.ID)
	if errMarkFeed != nil {
		log.Println("couldn't mark fetched feed:", errMarkFeed)
		return
	}

	// Get feed by URL
	feed, errGetFeedByURL := fetchFeed(context.Background(), nextFeed.Url)
	if errGetFeedByURL != nil {
		log.Println("couldn't fetch feed by URL:", errGetFeedByURL)
		return
	}

	// Uncomment this only if you want to see the feed's posts while scrapping	
	// printRSSFeed(feed)

	// Save posts to database
	for _, item := range feed.Channel.Item {
		publishedAt := sql.NullTime{}
		if parsedTime, errParseTime := time.Parse(time.RFC1123Z, item.PubDate); errParseTime == nil {
			publishedAt = sql.NullTime{
				Time: parsedTime,
				Valid: true,
			}
		} else {
			log.Println("couldn't parse time:", errParseTime)
		}
		_, errCreatePost := dbQueries.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Url:         item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid: true,
			},
			PublishedAt: publishedAt,
			FeedID:      nextFeed.ID,
		})
		if errCreatePost != nil && errCreatePost.Error() != `pq: duplicate key value violates unique constraint "posts_url_key"` {
			log.Println("couldn't insert post into database:", errCreatePost)
		}
	}
}

// Uncomment this only if you want to see the feed's posts while scrapping	
/* func printRSSFeed(feed *RSSFeed) {

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
		fmt.Printf("* Item %v:\n\n", i + 1)
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
} */