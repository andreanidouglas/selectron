package database

import "database/sql"

//Database interface for connection
type Database interface {
	Close()
	SQLCommandExec() (*sql.Rows, error)
}
