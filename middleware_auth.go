package main

import (
	"fmt"
	"net/http"

	database "github.com/uloamaka/rss_aggregator/internal/database"
	"github.com/uloamaka/rss_aggregator/internal/database/auth"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetApiKey(r.Header)
		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("Auth error: %v", err))
			return
		}
		user, err := apiCfg.DB.GetUserByApiKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Error geting user: %v", err))
			return
		}
		handler(w, r, user)
	}
} 