package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"net/url"
)

func DbConnect() (*sqlx.DB, error) {
	q := url.Values{}
	q.Set("sslmode", "disable")
	q.Set("timezone", "utc")

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword("dossier", "qwe"),
		Host:     "localhost:5432", // change here
		Path:     "dossier",
		RawQuery: q.Encode(),
	}
	fmt.Println("compiles creds", u.String())
	return sqlx.Open("postgres", u.String())
}
