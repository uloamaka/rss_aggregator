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
