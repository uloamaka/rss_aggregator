package main

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	database "github.com/uloamaka/rss_aggregator/internal/database"
)

func startScraping(
	db *database.Queries,
	concurrency int,
	timeBetweenRequest time.Duration,
) {
	log.Printf("Scraping on %v goroutines every %s duration", concurrency, timeBetweenRequest)
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
			go scrapeFeed(db, wg, feed)
		}
		wg.Wait()
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("Error marking feed as fetched:", err)
		return
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("Error fetching feed:", err)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		description := pgtype.Text{
			String: item.Description,
			Valid:  item.Description != "",
		}

		pubAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			// Try alternative date formats
			layouts := []string{
				time.RFC1123,
				time.RFC822Z,
				time.RFC822,
				"Mon, 02 Jan 2006 15:04:05 -0700",
				"2006-01-02T15:04:05Z",
			}
			
			parsed := false
			for _, layout := range layouts {
				if pubAt, err = time.Parse(layout, item.PubDate); err == nil {
					parsed = true
					break
				}
			}
			
			if !parsed {
				log.Printf("couldn't parse date %v with err %v", item.PubDate, err)
				continue
			}
		}

		now := time.Now()
		postID, err := uuid.NewRandom()
		if err != nil {
			log.Printf("error generating UUID: %v", err)
			continue
		}

		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID: pgtype.UUID{
				Bytes: postID,
				Valid: true,
			},
			Title: item.Title,
			Description: description,
			PublishedAt: pgtype.Timestamp{
				Time:  pubAt,
				Valid: true,
			},
			Url:    item.Link,
			FeedID: feed.ID,
			CreatedAt: pgtype.Timestamp{
				Time:  now,
				Valid: true,
			},
			UpdatedAt: pgtype.Timestamp{
				Time:  now,
				Valid: true,
			},
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Printf("Failed to create post %q: %v", item.Title, err)
			continue
		}
		log.Printf("Created post %q in feed %q", item.Title, feed.Name)
	}
	log.Printf("Feed %q collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}