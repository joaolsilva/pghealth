package postgres

import (
	"database/sql"
	_ "github.com/lib/pq"
	"time"
)

type PG struct {
	db *sql.DB
}

func newDatabaseConnection() (db *sql.DB, err error) {
	db, err = sql.Open("postgres", "host=localhost port=5432 dbname=postgres sslmode=disable")
	if err != nil {
		return db, err

	}
	db.SetConnMaxLifetime(10 * time.Minute)
	_, err = db.Exec("SET default_transaction_read_only = on;")
	return db, err
}

func NewPG() (*PG, error) {
	pg := PG{}
	db, err := newDatabaseConnection()
	pg.db = db
	return &pg, err
}

type Database struct {
	Name          string
	Size          int
	FormattedSize string
}

func (pg *PG) ListDatabases() (databases []Database, err error) {
	databases = []Database{}
	rows, err := pg.db.Query("SELECT datname,pg_database_size(datname) AS database_size, pg_size_pretty(pg_database_size(datname)) AS formatted_size FROM pg_stat_database WHERE datname IS NOT NULL ORDER BY database_size DESC;")
	if err != nil {
		return databases, err
	}
	defer rows.Close()
	for rows.Next() {
		d := Database{}
		err = rows.Scan(&d.Name, &d.Size, &d.FormattedSize)
		if err != nil {
			return databases, err
		}
		databases = append(databases, d)
	}

	return databases, nil
}
