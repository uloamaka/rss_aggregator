package main

import (
	"github.com/jackc/pgx/v5/pgtype"
	database "github.com/uloamaka/rss_aggregator/internal/database"
)

type User struct {
	ID        pgtype.UUID `json:"id"`
	Name      string `json:"name"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
	ApiKey string `json:"api_key"`
}

func databaseUserToUser(dbUser database.User) User {
	return User{
		ID: dbUser.ID,
		Name: dbUser.Name,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		ApiKey: dbUser.ApiKey,
	}
}

type Feed struct {
	ID        pgtype.UUID `json:"id"`
	Name      string `json:"name"`
	Url       string `json:"url"`
	UserID    pgtype.UUID `json:"user_id"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
}

func databaseFeedToFeed(dbFeed database.Feed) Feed {
	return Feed{
		ID: dbFeed.ID,
		Name: dbFeed.Name,
		Url: dbFeed.Url,
		UserID: dbFeed.UserID,
		CreatedAt: dbFeed.CreatedAt,
		UpdatedAt: dbFeed.UpdatedAt,
	}
}

func databaseFeedsToFeeds(dbFeeds []database.Feed) []Feed {
	feeds := []Feed{}
	for _, dbFeed := range dbFeeds{
		feeds = append(feeds, databaseFeedToFeed(dbFeed))
	}
	return feeds
}

type FeedFollow struct {
	ID        pgtype.UUID `json:"id"`
	UserID    pgtype.UUID `json:"user_id"`
	FeedID    pgtype.UUID `json:"feed_id"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
}

func databaseFeedFollowToFeedFollow(dbFeedFollow database.FeedFollow) FeedFollow {
	return FeedFollow{
		ID: dbFeedFollow.ID,
		UserID: dbFeedFollow.UserID,
		FeedID: dbFeedFollow.FeedID,
		CreatedAt: dbFeedFollow.CreatedAt,
		UpdatedAt: dbFeedFollow.UpdatedAt,
	}
}

func databaseFeedFollowsToFeedFollows(dbFeedFollows []database.FeedFollow) []FeedFollow {
	feedFollows := []FeedFollow{}
	for _, dbFeedFollow := range dbFeedFollows{
		feedFollows = append(feedFollows, databaseFeedFollowToFeedFollow(dbFeedFollow))
	}
	return feedFollows
}