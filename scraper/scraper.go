package scraper

// running in background to fetch data from RSS feed

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/LostArrows27/go-rss-aggregator/internal/database"
	"github.com/google/uuid"
)

/*
	input:

1. database connection
2. number of concurrent requests
3. time between requests
*/
func StartScrapping(
	db *database.Queries,
	concurrency int,
	timeBetweenRequsets time.Duration,
) {
	// 1. set up to run every timeBetweenRequsets
	log.Printf("Scrapping on %v gorooutines every %s duration", concurrency, timeBetweenRequsets)
	ticker := time.NewTicker(timeBetweenRequsets)
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(
			context.Background(),
			int32(concurrency),
		)

		if err != nil {
			log.Printf("Error fetching feeds: %v", err)
			continue
		}

		wg := &sync.WaitGroup{}

		for _, feed := range feeds {
			wg.Add(1)

			go scrapeFeed(db, wg, feed)
		}
		wg.Wait()

	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)

	if err != nil {
		log.Printf("Error marking feed as fetched: %v", err)
		return
	}

	rssFeed, err := URLToFeed(feed.Url)

	if err != nil {
		log.Printf("Error fetching feed: %v", err)
		return
	}

	for _, item := range rssFeed.Channel.Items {
		description := sql.NullString{}

		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}

		pubAt, err := time.Parse(time.RFC1123Z, item.PubDate)

		if err != nil {
			log.Printf("Error parsing time: %v", err)
			continue
		}

		_, err = db.CreatePost(
			context.Background(),
			database.CreatePostParams{
				ID:          uuid.New(),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				Title:       item.Title,
				Description: description,
				PublishedAt: pubAt,
				Url:         item.Link,
				FeedID:      feed.ID,
			},
		)

		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Printf("Error creating post: %v", err)
		}
	}

	log.Printf("Feed %s collected %v posts found", feed.Name, len(rssFeed.Channel.Items))
}
