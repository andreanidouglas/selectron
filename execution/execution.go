package execution

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

//Execution stores the parameters for execute a sql query
type Execution struct {
	Database  string
	Host      string
	Name      string
	Login     string
	Password  string
	SQLPath   string
	SQLString string
}

//New creates a new slice of execution parsing the contents of the given file
func New(executeFilePath string) ([]Execution, error) {
	f, err := os.Open(executeFilePath)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not open file path %s", executeFilePath))
	}

	l, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, errors.Wrap(err, "could not parse csv file")
	}

	e := make([]Execution, len(l))

	for i, line := range l {
		e[i].Database = line[0]
		e[i].Host = line[1]
		e[i].Name = line[2]
		e[i].Login = line[3]
		e[i].Password = line[4]
		e[i].SQLPath = line[5]
		err = e[i].readSQLString()
		if err != nil {
			return nil, fmt.Errorf("could not read sql file: %v", err)
		}
	}

	return e, nil
}

func (e *Execution) readSQLString() error {
	bytes, err := ioutil.ReadFile(e.SQLPath)
	if err != nil {
		return err
	}

	e.SQLString = string(bytes)
	return nil
}

// WriteResult writes a csv files with the result of the SQL query
func (e *Execution) WriteResult(rows *sql.Rows) error {

	filename := fmt.Sprintf("%s.csv", e.Name)
	f, err := os.Create("c:\\temp\\" + filename)
	if err != nil {
		return err
	}

	if rows.Err() != nil {
		return fmt.Errorf("error on rows: %v", rows.Err())
	}

	w := csv.NewWriter(f)
	headers, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("could get headers from SQL, %v", err)
	}

	w.Write(headers)
	w.Flush()

	dataRows := make([]interface{}, len(headers))

	for i := range headers {
		dataRows[i] = new(string)
	}

	for rows.Next() {

		err := rows.Scan(dataRows...)
		if err != nil {
			return fmt.Errorf("could not scan sql.Rows, %v", err)
		}
		row := make([]string, 0)
		for _, dataRow := range dataRows {
			s, ok := dataRow.(*string)
			if !ok {
				return fmt.Errorf("could not assert type")
			}
			row = append(row, *s)
		}

		err = w.Write(row)
		if err != nil {
			return fmt.Errorf("could not write results to file, %v", err)
		}
	}

	w.Flush()
	f.Close()

	return nil
}
