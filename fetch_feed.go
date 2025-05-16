package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"
)

// Save channel info
type RSSFeed struct {
	Channel struct {
		Title string `xml:"title"`
		Link  []struct {
			Text string `xml:",chardata"`
			Href string `xml:"href,attr"`
			Rel  string `xml:"rel,attr"`
			Type string `xml:"type,attr"`
		} `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

// Save item info that is within channel
type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

// Fetch feed from a given URL
func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {

	// Create the new request
	req, errNewReq := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if errNewReq != nil {
		return &RSSFeed{}, fmt.Errorf("couldn't create new request with context: %w", errNewReq)
	}

	// Set User-Agent header and initialize client
	req.Header.Set("User-Agent", "gator")

	client := &http.Client{Timeout: time.Second * 5}

	// Do the request through the client
	resp, errDo := client.Do(req)
	if errDo != nil {
		return &RSSFeed{}, fmt.Errorf("couldn't do the request: %w", errDo)
	}

	defer resp.Body.Close()

	// Read response body
	bodyData, errRead := io.ReadAll(resp.Body)
	if errRead != nil {
		return &RSSFeed{}, fmt.Errorf("couldn't read response body: %w", errRead)
	}

	// Unmarshal response body
	var rssFeed RSSFeed
	errUnmarshal := xml.Unmarshal(bodyData, &rssFeed)
	if errUnmarshal != nil {
		return &RSSFeed{}, fmt.Errorf("couldn't unmarshal response body: %w", errUnmarshal)
	}

	// Unescape Titles and Descriptions
	unescapeRSS(&rssFeed)

	return &rssFeed, nil
}

// Unescape Titles and Descriptions using html.Unescapestring()
func unescapeRSS(feed *RSSFeed) {

	// Change of fields
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i := 0; i < len(feed.Channel.Item); i++ {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
	}
}
