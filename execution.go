package main

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/pkg/errors"
)

//Execution stores the parameters for execute a sql query
type Execution struct {
	Database string
	Host     string
	Name     string
	Login    string
	Password string
	SQLPath  string
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
	}

	return e, nil
}
