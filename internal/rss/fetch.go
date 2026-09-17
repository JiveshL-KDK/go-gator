package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"net/http"
)

func decodeStrings(feed *RSSFeed) {
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
}

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)

	if err != nil {
		return nil, fmt.Errorf("Can't make req: %w", err)
	}

	req.Header.Set("User-Agent", "gator")

	client := http.Client{}

	res, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("Can't do req: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unsucessful req: %w", err)
	}

	var feed RSSFeed
	if err := xml.NewDecoder(res.Body).Decode(&feed); err != nil {
		return nil, fmt.Errorf("Failed to parse res: %w", err)
	}

	decodeStrings(&feed)

	return &feed, nil

}
