package db

import (
	"github.com/mager/bluedot/config"

	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host   = "ep-black-dream-58840312.us-east-2.aws.neon.tech"
	port   = 5432
	user   = "mager"
	dbname = "neondb"
)

func ProvideDB(
	cfg config.Config,
) *sql.DB {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", host, port, user, cfg.PGPassword, dbname)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	// check db
	err = db.Ping()
	CheckError(err)

	// Log message including hostname
	fmt.Println("Connected to database on", host, "as", user)

	return db
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

var Options = ProvideDB
