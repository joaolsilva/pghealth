package postgres

type TableCacheHitRatio struct {
	Table            string `db:"table_name"`
	DiskHits         int    `db:"disk_hits" table:"Disk Hits"`
	PercentDiskHits  string `db:"percent_disk_hits" table:"% Disk Hits"`
	PercentCacheHits string `db:"percent_cache_hits" table:"% Cache Hits"`
	TotalHits        int    `db:"total_hits" table:"Total Hits"`
}

func (dbConnection *DBConnection) GetTableCacheHitRatios() (tableCacheHitRatios []TableCacheHitRatio, err error) {
	tableCacheHitRatios = []TableCacheHitRatio{}
	err = dbConnection.db.Select(&tableCacheHitRatios, `

WITH
all_tables as
(
SELECT  *
FROM    (
    SELECT  'all'::text as table_name,
        sum( (coalesce(heap_blks_read,0) + coalesce(idx_blks_read,0) + coalesce(toast_blks_read,0) + coalesce(tidx_blks_read,0)) ) as from_disk,
        sum( (coalesce(heap_blks_hit,0)  + coalesce(idx_blks_hit,0)  + coalesce(toast_blks_hit,0)  + coalesce(tidx_blks_hit,0))  ) as from_cache
    FROM    pg_statio_all_tables  --> change to pg_statio_USER_tables if you want to check only user tables (excluding postgres's own tables)
    ) a
WHERE   (from_disk + from_cache) > 0 -- discard tables without hits
),
tables as
(
SELECT  *
FROM    (
    SELECT  relname as table_name,
        ( (coalesce(heap_blks_read,0) + coalesce(idx_blks_read,0) + coalesce(toast_blks_read,0) + coalesce(tidx_blks_read,0)) ) as from_disk,
        ( (coalesce(heap_blks_hit,0)  + coalesce(idx_blks_hit,0)  + coalesce(toast_blks_hit,0)  + coalesce(tidx_blks_hit,0))  ) as from_cache
    FROM    pg_statio_all_tables --> change to pg_statio_USER_tables if you want to check only user tables (excluding postgres's own tables)
    ) a
WHERE   (from_disk + from_cache) > 0 -- discard tables without hits
)
SELECT  table_name,
    from_disk as disk_hits,
    round((from_disk::numeric / (from_disk + from_cache)::numeric)*100.0,2) as percent_disk_hits,
    round((from_cache::numeric / (from_disk + from_cache)::numeric)*100.0,2) as percent_cache_hits,
    (from_disk + from_cache) as total_hits
FROM    (SELECT * FROM all_tables UNION ALL SELECT * FROM tables) a
ORDER   BY (case when table_name = 'all' then 0 else 1 end), from_disk desc;
`)
	return tableCacheHitRatios, err
}
