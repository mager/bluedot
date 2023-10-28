package db

import (
	"database/sql"
)

type User struct {
	// Order matters here
	ID            string
	Name          string
	Email         string
	EmailVerified sql.NullTime
	Image         string
	Slug          string
}

type Dataset struct {
	// Order matters here
	ID          string
	UserID      string
	Name        string
	Slug        string
	Source      string
	Description sql.NullString
	Created     sql.NullTime
	Updated     sql.NullTime
}
