package postgres

type VacuumStat struct {
	Table               string `db:"table"`
	LastVacuum          string `db:"last_vacuum" table:"Last Vacuum"`
	LastAutoVacuum      string `db:"last_autovacuum" table:"Last Auto-Vacuum"`
	RowCount            string `db:"rowcount" table:"Row Count"`
	DeadRowCount        string `db:"dead_rowcount" table:"Dead Row Count"`
	RowsPerPage         string `db:"rows_per_page" table:"Rows per Page"`
	AutoVacuumThreshold string `db:"autovacuum_threshold" table:"Auto-Vacuum Threshold"`
	WillVacuum          string `db:"will_vacuum" table:"Will Vacuum"`
}

func (dbConnection *DBConnection) GetVacuumStats() (vacuumStats []VacuumStat, err error) {
	vacuumStats = []VacuumStat{}
	err = dbConnection.db.Select(&vacuumStats, `
WITH table_opts AS (
  SELECT
    pg_class.oid, relname, nspname, array_to_string(reloptions, '') AS relopts
  FROM
     pg_class INNER JOIN pg_namespace ns ON relnamespace = ns.oid
), vacuum_settings AS (
  SELECT
    oid, relname, nspname,
    CASE
      WHEN relopts LIKE '%autovacuum_vacuum_threshold%'
        THEN substring(relopts, '.*autovacuum_vacuum_threshold=([0-9.]+).*')::integer
        ELSE current_setting('autovacuum_vacuum_threshold')::integer
      END AS autovacuum_vacuum_threshold,
    CASE
      WHEN relopts LIKE '%autovacuum_vacuum_scale_factor%'
        THEN substring(relopts, '.*autovacuum_vacuum_scale_factor=([0-9.]+).*')::real
        ELSE current_setting('autovacuum_vacuum_scale_factor')::real
      END AS autovacuum_vacuum_scale_factor
  FROM
    table_opts
)
SELECT
  vacuum_settings.relname AS table,
  COALESCE(to_char(psut.last_vacuum, 'YYYY-MM-DD HH24:MI'), '') AS last_vacuum,
  COALESCE(to_char(psut.last_autovacuum, 'YYYY-MM-DD HH24:MI'), '') AS last_autovacuum,
  COALESCE(to_char(pg_class.reltuples, '9G999G999G999'), '') AS rowcount,
  COALESCE(to_char(psut.n_dead_tup, '9G999G999G999'), '') AS dead_rowcount,
  COALESCE(to_char(pg_class.reltuples / NULLIF(pg_class.relpages, 0), '999G999.99'), '') AS rows_per_page,
  COALESCE(to_char(autovacuum_vacuum_threshold
       + (autovacuum_vacuum_scale_factor::numeric * pg_class.reltuples), '9G999G999G999'), '') AS autovacuum_threshold,
  CASE
    WHEN autovacuum_vacuum_threshold + (autovacuum_vacuum_scale_factor::numeric * pg_class.reltuples) < psut.n_dead_tup
    THEN 'yes' ELSE ''
  END AS will_vacuum
FROM
  pg_stat_user_tables psut INNER JOIN pg_class ON psut.relid = pg_class.oid
    INNER JOIN vacuum_settings ON pg_class.oid = vacuum_settings.oid
ORDER BY psut.n_dead_tup DESC, pg_class.reltuples DESC;
`)
	return vacuumStats, err
}
