package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	database "github.com/uloamaka/rss_aggregator/internal/database"
)

func (apiCfg *apiConfig)handlerCreateFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID pgtype.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing json: %v", err))
		return  
	}
	
	userID, err := uuid.NewRandom()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error generating UUID: %v", err))
		return
	}
	now := time.Now()
	feedFollow, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams {
		ID:        pgtype.UUID{Bytes: [16]byte(userID), Valid: true},
		UserID: user.ID,
		FeedID: params.FeedID,
		CreatedAt: pgtype.Timestamp{Time: now, Valid: true},
		UpdatedAt: pgtype.Timestamp{Time: now, Valid: true},
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error creating feed follow: %v", err))
		return
	}
	respondWithJson(w, 201, databaseFeedFollowToFeedFollow(feedFollow))
}

func (apiCfg *apiConfig)handlerGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	
	feedFollows, err := apiCfg.DB.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error fetching feed follows: %v", err))
		return
	}
	respondWithJson(w, 201, databaseFeedFollowsToFeedFollows(feedFollows))
}