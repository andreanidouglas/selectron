package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	oci8 "github.com/mattn/go-oci8"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage ./%s <filename>", os.Args[0])
	}
	var execs []Execution

	execs, err := New(os.Args[1])
	if err != nil {
		log.Fatalf("could not start execution: %v", err)
	}

	for _, exec := range execs {
		go func(exec Execution) {

		}(exec)
	}

}

func buildConnectionString(e Execution) string {
	return fmt.Sprintf("%s/%s@%s", e.Login, e.Password)
}

func database() {

	oci8.OCI8Driver.Logger = log.New(os.Stderr, "oci8 ", log.Ldate|log.Ltime|log.LUTC|log.Llongfile)

	connString := "andredr/D5022a38@p601.noa.alcoa.com"

	db, err := sql.Open("oci8", connString)
	if err != nil {
		log.Fatalf("could not open database: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 55*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, "select 1 from dual")
	if err != nil {
		log.Fatalf("could not query context: %v", err)
	}

	if !rows.Next() {
		log.Fatalf("no Next rows")

	}

	dest := make([]interface{}, 1)
	destPointer := make([]interface{}, 1)

	destPointer[0] = &dest[0]
	err = rows.Scan(destPointer...)

	if err != nil {
		log.Fatalf("could not read results: %v", err)
	}

	data, ok := dest[0].(float64)

	if !ok {
		log.Fatalf("could not smash results")
	}
	fmt.Println(data)

	cancel()

}
