package database

import "database/sql"

//Database interface for connection
type Database interface {
	BuildConnectionString(login string, password string, host string) string
	Close()
	SQLCommandExec(sqlCommand string, params ...string) (*sql.Rows, error)
}
