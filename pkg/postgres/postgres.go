package postgres

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type PG struct {
	db *sql.DB
}

func NewPG() (*PG, error) {
	pg := PG{}
	db, err := connectToDB(defaultDatabase)
	pg.db = db
	return &pg, err
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
