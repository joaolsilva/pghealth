package postgres

type TableSize struct {
	Schema      string `db:"table_schema"`
	TableName   string `db:"table_name" table:"Table Name"`
	RowEstimate string `db:"row_estimate" table:"Row Estimate"`
	Table       string `db:"table"`
	Index       string `db:"index"`
	TOAST       string `db:"toast"`
	Total       string `db:"total"`
}

func (dbConnection *DBConnection) GetTableSizes() (tableSizes []TableSize, err error) {
	tableSizes = []TableSize{}
	err = dbConnection.db.Select(&tableSizes, `SELECT
       table_schema,
       table_name,
       row_estimate,
       COALESCE(pg_size_pretty(table_bytes), '') AS table,
       COALESCE(pg_size_pretty(index_bytes), '') AS index,
       COALESCE(pg_size_pretty(toast_bytes), '') AS toast,
       COALESCE(pg_size_pretty(total_bytes), '') AS total
FROM
     (SELECT *, total_bytes-index_bytes-coalesce(toast_bytes,0) AS table_bytes
     FROM ( SELECT nspname AS table_schema, relname AS table_name , c.reltuples AS row_estimate , pg_total_relation_size(c.oid) AS total_bytes , pg_indexes_size(c.oid) AS index_bytes , pg_total_relation_size(reltoastrelid) AS toast_bytes FROM pg_class c LEFT JOIN pg_namespace n ON n.oid = c.relnamespace WHERE relkind = 'r' ) a ) a
ORDER BY total_bytes DESC;
`)
	return tableSizes, err
}
