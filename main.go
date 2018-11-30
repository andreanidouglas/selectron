package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/andreanidouglas/selectron/database/oracle"
	"github.com/andreanidouglas/selectron/execution"
)

type fileLog struct {
	f  *os.File
	RW sync.RWMutex
}

func main() {
	var export bool
	var outputLog string
	var inputFile string
	flag.BoolVar(&export, "export", false, "exports the template file")
	flag.StringVar(&outputLog, "log", "", "select where to export log")
	flag.StringVar(&inputFile, "execute", "export.csv", "define file to execute")
	flag.Parse()

	if export {
		f, err := os.Create("export.csv")
		if err != nil {
			log.Fatalf("could not create export: %v", err)
		}
		w := csv.NewWriter(f)

		header := []string{"database type", "host", "report name", "login", "password", "sql path"}
		err = w.Write(header)
		if err != nil {
			log.Fatalf("could not write export: %v", err)
		}
		w.Flush()
		f.Close()

		os.Exit(0)
	}

	f := new(fileLog)

	f.f = os.Stderr
	if outputLog != "" {
		var err error
		f.f, err = os.Create(outputLog)
		if err != nil {
			log.Fatalf("could not create error: %v", err)
		}
	}

	var wg sync.WaitGroup

	var execs []execution.Execution

	execs, err := execution.New(inputFile)

	if err != nil {
		log.Fatalf("could not start execution: %v", err)
	}

	wg.Add(len(execs))
	for _, exec := range execs {
		go func(ex execution.Execution) {
			f.RW.Lock()
			fmt.Fprintf(f.f, "\nStarting: %s", ex.Name)
			f.RW.Unlock()
			db, err := oracle.New(ex)
			if err != nil {
				log.Fatalf("could not create new db connection: %v", err)
			}
			rows, err := db.SQLCommandExec()
			if err != nil {
				log.Fatalf("could not execute sql command: %v", err)
			}
			err = ex.WriteResult(rows)
			if err != nil {
				log.Fatalf("could not write result to file: %v", err)
			}

			db.Close()
			f.RW.Lock()
			fmt.Fprintf(f.f, "\nCompleted : %s", ex.Name)
			f.RW.Unlock()
			wg.Done()
		}(exec)
	}

	wg.Wait()

	f.RW.Lock()
	fmt.Fprint(f.f, "\nAll Tasks Completed")
	f.RW.Unlock()
}
