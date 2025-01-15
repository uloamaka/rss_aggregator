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

func (apiCfg *apiConfig)handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
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
	user, err := apiCfg.DB.Createuser(r.Context(), database.CreateuserParams {
		ID:        pgtype.UUID{Bytes: [16]byte(userID), Valid: true},
		Name:      params.Name,
		CreatedAt: pgtype.Timestamp{Time: now, Valid: true},
		UpdatedAt: pgtype.Timestamp{Time: now, Valid: true},
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error creating user: %v", err))
		return
	}
	respondWithJson(w, 200, user)
}