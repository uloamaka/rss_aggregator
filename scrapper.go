package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	database "github.com/uloamaka/rss_aggregator/internal/database"
)
func startScraping(
	db *database.Queries,
	concurrency int,
	timeBetweenRequest time.Duration,
) {
	fmt.Printf("Scraping on %v goroutine every %s duration", concurrency, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest)
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(
			context.Background(),
			int32(concurrency),
		 )
		 if err != nil {
			log.Println("error fetching feeds:", err)
			continue
		 }
		 wg := &sync.WaitGroup{}
		 for _, feed := range feeds {
			wg.Add(1)

			go scrapFeed(db, wg, feed)
		 }
		 wg.Wait()
	}
}

func scrapFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("Error marking feed as fetched", err)
		return
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("Error fetching feed:", err)
		return
	}
	for _, item := range rssFeed.Channel.Item {
		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}

		pubAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Println("couldn't parse date %v with err %v", item.PubDate, err)
			continue
		}
		db.CreatePost(context.Background(), database.CreatePostParams{
			ID: pgtype.UUID{},
			Title: item.Title,
			Description: description,
			PublishedAt: pubAt,
			Url: item.Link,
			FeedID: feed.ID,
			CreatedAt: pgtype.Timestamp{Time: now, Valid: true},
			UpdatedAt: pgtype.Timestamp{Time: now, Valid: true},
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Println("Failed to create posts:",err)
		}
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}
