package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

type DatabaseName string

const defaultDatabase DatabaseName = "postgres"

type DBConnection struct {
	db *sql.DB
}

type Database struct {
	Name           DatabaseName
	Size           int
	FormattedSize  string
	CommitRatio    string
	CacheHitRation string
	BlocksRead     int
}

func connectToDB(dbName DatabaseName) (db *sql.DB, err error) {
	db, err = sql.Open("postgres", fmt.Sprintf("host=localhost port=5432 dbname=%v sslmode=disable", dbName))
	if err != nil {
		return db, err
	}
	db.SetConnMaxLifetime(10 * time.Minute)
	_, err = db.Exec("SET default_transaction_read_only = on;")
	return db, err
}

func NewDBConnection(dbName DatabaseName) (*DBConnection, error) {
	dbConnection := DBConnection{}
	db, err := connectToDB(dbName)
	dbConnection.db = db
	return &dbConnection, err
}

type Bloat struct {
	Type       string
	SchemaName string
	ObjectName string
	Bloat      string
	Waste      string
}

func (dbConnection *DBConnection) ListBloat() (bloat []Bloat, err error) {
	bloat = []Bloat{}
	rows, err := dbConnection.db.Query(`
WITH constants AS (
  SELECT current_setting('block_size')::numeric AS bs, 23 AS hdr, 4 AS ma
), bloat_info AS (
  SELECT
    ma,bs,schemaname,tablename,
    (datawidth+(hdr+ma-(case when hdr%ma=0 THEN ma ELSE hdr%ma END)))::numeric AS datahdr,
    (maxfracsum*(nullhdr+ma-(case when nullhdr%ma=0 THEN ma ELSE nullhdr%ma END))) AS nullhdr2
  FROM (
    SELECT
      schemaname, tablename, hdr, ma, bs,
      SUM((1-null_frac)*avg_width) AS datawidth,
      MAX(null_frac) AS maxfracsum,
      hdr+(
        SELECT 1+count(*)/8
        FROM pg_stats s2
        WHERE null_frac<>0 AND s2.schemaname = s.schemaname AND s2.tablename = s.tablename
      ) AS nullhdr
    FROM pg_stats s, constants
    GROUP BY 1,2,3,4,5
  ) AS foo
), table_bloat AS (
  SELECT
    schemaname, tablename, cc.relpages, bs,
    CEIL((cc.reltuples*((datahdr+ma-
      (CASE WHEN datahdr%ma=0 THEN ma ELSE datahdr%ma END))+nullhdr2+4))/(bs-20::float)) AS otta
  FROM bloat_info
  JOIN pg_class cc ON cc.relname = bloat_info.tablename
  JOIN pg_namespace nn ON cc.relnamespace = nn.oid AND nn.nspname = bloat_info.schemaname AND nn.nspname <> 'information_schema'
), index_bloat AS (
  SELECT
    schemaname, tablename, bs,
    COALESCE(c2.relname,'?') AS iname, COALESCE(c2.reltuples,0) AS ituples, COALESCE(c2.relpages,0) AS ipages,
    COALESCE(CEIL((c2.reltuples*(datahdr-12))/(bs-20::float)),0) AS iotta -- very rough approximation, assumes all cols
  FROM bloat_info
  JOIN pg_class cc ON cc.relname = bloat_info.tablename
  JOIN pg_namespace nn ON cc.relnamespace = nn.oid AND nn.nspname = bloat_info.schemaname AND nn.nspname <> 'information_schema'
  JOIN pg_index i ON indrelid = cc.oid
  JOIN pg_class c2 ON c2.oid = i.indexrelid
)
SELECT
  type, schemaname, object_name, bloat, pg_size_pretty(raw_waste) as waste
FROM
(SELECT
  'table' as type,
  schemaname,
  tablename as object_name,
  ROUND(CASE WHEN otta=0 THEN 0.0 ELSE table_bloat.relpages/otta::numeric END,1) AS bloat,
  CASE WHEN relpages < otta THEN '0' ELSE (bs*(table_bloat.relpages-otta)::bigint)::bigint END AS raw_waste
FROM
  table_bloat
    UNION
SELECT
  'index' as type,
  schemaname,
  tablename || '::' || iname as object_name,
  ROUND(CASE WHEN iotta=0 OR ipages=0 THEN 0.0 ELSE ipages/iotta::numeric END,1) AS bloat,
  CASE WHEN ipages < iotta THEN '0' ELSE (bs*(ipages-iotta))::bigint END AS raw_waste
FROM
  index_bloat) bloat_summary
ORDER BY raw_waste DESC, bloat DESC;
`)
	if err != nil {
		return bloat, err
	}
	defer rows.Close()
	for rows.Next() {
		b := Bloat{}
		err = rows.Scan(&b.Type, &b.SchemaName, &b.ObjectName, &b.Bloat, &b.Waste)
		if err != nil {
			return bloat, err
		}
		bloat = append(bloat, b)
	}

	return bloat, nil
}
