package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"Blog-aggregator/internal/database"

	"github.com/google/uuid"
)

func scrapeFeeds(s *state) error {
	ctx := context.Background()
	nFeed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return err
	}

	paramsMFF := database.MarkFeedFetchedParams{
		ID:            nFeed.ID,
		LastFetchedAt: nullTime(time.Now()),
	}
	if err := s.db.MarkFeedFetched(ctx, paramsMFF); err != nil {
		return err
	}

	rss, err := fetchFeed(ctx, nFeed.Url)
	if err != nil {
		return err
	}

	fmt.Printf("\n\nFeed: %v\n\n", nFeed.Name)
	rss.Channel.Description = html.UnescapeString(rss.Channel.Description)
	rss.Channel.Title = html.UnescapeString(rss.Channel.Title)

	const timeFormat = "Mon, 02 Jan 2006 15:04:05 -0700"

	for _, item := range rss.Channel.Item {
		title, desc := item.Title, item.Description
		item.Description = html.UnescapeString(desc)
		item.Title = html.UnescapeString(title)
		pubDate, err := time.Parse(timeFormat, item.PubDate)
		if err != nil {
			return err
		}

		paramsCP := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   nullTime(time.Now()),
			Title:       item.Title,
			Description: item.Description,
			Url:         item.Link,
			PublishedAt: pubDate,
			FeedID:      nFeed.ID,
		}

		if err := s.db.CreatePost(ctx, paramsCP); err != nil {
			if err.Error() == `pq: duplicate key value violates unique constraint "posts_url_key"` {
				continue
			}
			return err
		}
	}

	return nil
}

func nullTime(t time.Time) sql.NullTime {
	return sql.NullTime{
		Time:  t,
		Valid: true,
	}
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	f := func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.UserName)
		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}
	return f
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}

	var feed RSSFeed
	if err := xml.Unmarshal(body, &feed); err != nil {
		return nil, err
	}

	return &feed, nil
}
