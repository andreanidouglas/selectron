package oracle

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/andreanidouglas/selectron/database"
	"github.com/andreanidouglas/selectron/execution"
	oci8 "github.com/mattn/go-oci8"
	"github.com/pkg/errors"
)

type oracle struct {
	dbc        *sql.DB
	sqlCommand string
}

// New creates a new oracle database connection
func New(e execution.Execution) (database.Database, error) {
	oci8.OCI8Driver.Logger = log.New(os.Stderr, "oci8 ", log.Ldate|log.Ltime|log.LUTC|log.Llongfile)

	connString := buildConnectionString(e.Login, e.Password, e.Host)

	dbc, err := sql.Open("oci8", connString)
	if err != nil {
		log.Fatalf("could not open database: %v", err)
	}
	return &oracle{dbc, e.SQLString}, nil
}

func buildConnectionString(login string, password string, host string) string {

	return fmt.Sprintf("%s/%s@%s", login, password, host)

}

func (db *oracle) SQLCommandExec() (*sql.Rows, error) {

	rows, err := db.dbc.Query(db.sqlCommand)
	if err != nil {
		return nil, errors.Wrap(err, "could not query on context")
	}

	return rows, nil
}

func (db *oracle) Close() {
	db.dbc.Close()
}
