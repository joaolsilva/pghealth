package postgres

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PG struct {
	db *sqlx.DB
}

func NewPG() (*PG, error) {
	pg := PG{}
	db, err := connectToDB(defaultDatabase)
	pg.db = db
	return &pg, err
}

func (pg *PG) ListDatabases() (databases []Database, err error) {
	databases = []Database{}

	err = pg.db.Select(&databases, "SELECT datname,pg_database_size(datname) AS size, pg_size_pretty(pg_database_size(datname)) AS formatted_size,CASE WHEN (xact_commit + xact_rollback)>0 THEN round(100.0 * xact_commit / (xact_commit + xact_rollback),1)::varchar ELSE '' END as commit_ratio,CASE WHEN (blks_hit + blks_read) > 0 THEN round(100.0 * blks_hit / (blks_hit + blks_read),1)::varchar ELSE '' END as cache_hit_ratio, blks_read FROM pg_stat_database WHERE datname IS NOT NULL ORDER BY size DESC;")
	return databases, err
}
