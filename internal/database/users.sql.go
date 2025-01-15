// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createuser = `-- name: Createuser :one
INSERT INTO users (id, name, created_at, updated_at)
VALUES($1, $2, $3, $4)
RETURNING id, name, created_at, updated_at
`

type CreateuserParams struct {
	ID        pgtype.UUID
	Name      string
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}

func (q *Queries) Createuser(ctx context.Context, arg CreateuserParams) (User, error) {
	row := q.db.QueryRow(ctx, createuser,
		arg.ID,
		arg.Name,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
