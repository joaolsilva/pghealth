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
	Name           string
	Size           int
	FormattedSize  string
	CommitRatio    string
	CacheHitRation string
	BlocksRead     int
}

func (pg *PG) ListDatabases() (databases []Database, err error) {
	databases = []Database{}
	rows, err := pg.db.Query("SELECT datname,pg_database_size(datname) AS database_size, pg_size_pretty(pg_database_size(datname)) AS formatted_size,CASE WHEN (xact_commit + xact_rollback)>0 THEN round(100.0 * xact_commit / (xact_commit + xact_rollback),1)::varchar ELSE '' END as commit_ratio,CASE WHEN (blks_hit + blks_read) > 0 THEN round(100.0 * blks_hit / (blks_hit + blks_read),1)::varchar ELSE '' END as cache_hit_ratio, blks_read FROM pg_stat_database WHERE datname IS NOT NULL ORDER BY database_size DESC;")
	if err != nil {
		return databases, err
	}
	defer rows.Close()
	for rows.Next() {
		d := Database{}
		err = rows.Scan(&d.Name, &d.Size, &d.FormattedSize, &d.CommitRatio, &d.CacheHitRation, &d.BlocksRead)
		if err != nil {
			return databases, err
		}
		databases = append(databases, d)
	}

	return databases, nil
}
