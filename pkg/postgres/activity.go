package postgres

type Activity struct {
	ProcessID  string `db:"pid" table:"Process ID"`
	QueryStart string `db:"query_start" table:"Query Start"`
	State      string `db:"state"`
	Query      string `db:"query"`
}

func (dbConnection *DBConnection) GetActivity() (activity []Activity, err error) {
	activity = []Activity{}
	err = dbConnection.db.Select(&activity, `
SELECT pid,query_start,state,query
FROM pg_stat_activity
WHERE datname=current_database() AND pid != pg_backend_pid()
ORDER BY query_start;
`)
	return activity, err
}
