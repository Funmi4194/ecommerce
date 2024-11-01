package database

import (
	"database/sql"
	"os"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

func NewPostgreSQLConnection(uri string, connections int32, debug bool) error {
	// creates a new connection to postgreSQL with the given uri
	pgDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(uri)))

	// sets the maximum number of open connections to the database.
	pgDB.SetMaxOpenConns(int(connections))

	// initializes a new Bun DB instance using the open `pgDB` connection.
	PostgreSQLDB = bun.NewDB(pgDB, pgdialect.New())

	// adds a query hook for debugging if `debug` is set to true.
	// this hook logs detailed query information for debugging purposes.
	if debug {
		PostgreSQLDB.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}

	return PostgreSQLDB.Ping()
}

func ReadFileAndExecuteQueries(path string) error {
	queries, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	_, err = PostgreSQLDB.Exec(string(queries))
	if err != nil {
		return err
	}
	return nil
}
