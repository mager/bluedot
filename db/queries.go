package db

import (
	"database/sql"
	"fmt"
)

func GetUserByUsername(db *sql.DB, username string) User {
	var user User
	row := db.QueryRow("SELECT id, name, email, image, slug FROM \"User\" WHERE slug=$1", username)
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Image,
		&user.Slug,
	)
	CheckError(err)
	return user
}

func GetDatasetByUserIdAndSlug(db *sql.DB, userID string, slug string) Dataset {
	var dataset Dataset
	queryString := fmt.Sprintf("SELECT id, \"userId\", name, slug, source, description, created, updated FROM \"Dataset\" WHERE \"userId\" = '%s' AND slug = '%s'", userID, slug)
	row := db.QueryRow(queryString)
	// Log the row to the console
	fmt.Println(row)

	err := row.Scan(
		&dataset.ID,
		&dataset.UserID,
		&dataset.Name,
		&dataset.Slug,
		&dataset.Source,
		&dataset.Description,
		&dataset.Created,
		&dataset.Updated,
	)
	CheckError(err)
	return dataset
}
