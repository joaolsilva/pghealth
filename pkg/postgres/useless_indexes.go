package postgres

type UselessIndexes struct {
	Table      string `db:"table"`
	Index      string `db:"index"`
	IndexSize  string `db:"index_size" table:"Index Size"`
	IndexScans int    `db:"index_scans" table:"Index Scans"`
}

func (dbConnection *DBConnection) GetUselessIndexes() (uselessIndexes []UselessIndexes, err error) {
	uselessIndexes = []UselessIndexes{}
	err = dbConnection.db.Select(&uselessIndexes, `
SELECT
       schemaname || '.' || relname AS table,
       indexrelname AS index,
       pg_size_pretty(pg_relation_size(i.indexrelid)) AS index_size,
       idx_scan as index_scans
FROM
     pg_stat_user_indexes ui JOIN pg_index i ON ui.indexrelid = i.indexrelid
WHERE
      NOT indisunique AND idx_scan < 50 AND pg_relation_size(relid) > 5 * 8192
ORDER BY
         pg_relation_size(i.indexrelid) / nullif(idx_scan, 0) DESC NULLS FIRST,
         pg_relation_size(i.indexrelid) DESC;
`)
	return uselessIndexes, err
}
