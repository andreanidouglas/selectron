package oracle

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/andreanidouglas/sql-go/database"
	oci8 "github.com/mattn/go-oci8"
	"github.com/pkg/errors"
)

type oracle struct {
	dbc *sql.DB
}

func New() (database.Database, error) {
	oci8.OCI8Driver.Logger = log.New(os.Stderr, "oci8 ", log.Ldate|log.Ltime|log.LUTC|log.Llongfile)

	connString := "andredr/D5022a38@p601.noa.alcoa.com"

	dbc, err := sql.Open("oci8", connString)
	if err != nil {
		log.Fatalf("could not open database: %v", err)
	}
	return &oracle{dbc}, nil
}

func (db *oracle) BuildConnectionString(login string, password string, host string) string {

	return fmt.Sprintf("%s/%s@%s", login, password, host)

}

func (db *oracle) Close() {
	db.dbc.Close()
}

func (db *oracle) SqlCommandExec(sqlCommand string, params ...string) (*sql.Rows, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	rows, err := db.dbc.QueryContext(ctx, sqlCommand, params)
	if err != nil {
		return nil, errors.Wrap(err, "could not query on context")
	}
	return rows, nil
}
