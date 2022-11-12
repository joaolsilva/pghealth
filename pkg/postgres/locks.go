package postgres

type ActiveLock struct {
	RelationName       string `db:"relname" table:"Relation Name"`
	LockType           string `db:"locktype" table:"Lock Type"`
	Page               string `db:"page"`
	VirtualTransaction string `db:"virtualtransaction" table:"Virtual Transaction"`
	PID                string `db:"pid"`
	Mode               string `db:"mode"`
	Granted            string `db:"granted"`
}

func (dbConnection *DBConnection) GetActiveLocks() (activeLocks []ActiveLock, err error) {
	activeLocks = []ActiveLock{}
	err = dbConnection.db.Select(&activeLocks, `
SELECT t.relname,
       COALESCE(l.locktype, '') AS locktype,
       COALESCE(l.page::varchar, '') AS page,
       COALESCE(l.virtualtransaction, '') AS virtualtransaction,
       COALESCE(l.pid::varchar, '') AS pid,
       COALESCE(l.mode, '') AS mode,
       COALESCE(l.granted::varchar, '') AS granted
FROM pg_locks l, pg_stat_all_tables t
WHERE l.relation=t.relid;
`)
	return activeLocks, err
}
