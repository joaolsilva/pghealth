package postgres

type MissingIndexes struct {
	RelationName    string `db:"relname"`
	TooMuchSeq      int    `db:"too_much_seq" table:"Too Much Seq"`
	MissingIndex    string `db:"missing_index" table:"Missing Index?"`
	RelationSize    string `db:"rel_size" table:"Size"`
	SequentialScans int    `db:"seq_scan" table:"Sequential Scans"`
	IndexScans      int    `db:"idx_scan" table:"Index Scans"`
}

func (dbConnection *DBConnection) GetMissingIndexes() (missingIndexes []MissingIndexes, err error) {
	missingIndexes = []MissingIndexes{}
	err = dbConnection.db.Select(&missingIndexes, `SELECT
       relname,
       seq_scan-idx_scan AS too_much_seq,
       case when seq_scan-idx_scan>0 THEN 'Missing Index?' ELSE 'OK' END AS missing_index,
       pg_size_pretty(pg_relation_size(relname::regclass)) AS rel_size,
       seq_scan,
       idx_scan
FROM
     pg_stat_all_tables
WHERE
      schemaname='public'
  AND pg_relation_size(relname::regclass)>80000
ORDER BY too_much_seq DESC;
`)
	return missingIndexes, err
}
