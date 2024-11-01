package main

import (
	"net/http"
	"os"

	"github.com/funmi4194/ecommerce/database"
	"github.com/funmi4194/ecommerce/database/migration"
	"github.com/funmi4194/ecommerce/primer"
	"github.com/funmi4194/ecommerce/version"
	"github.com/opensaucerer/barf"
)

func main() {
	// set env path
	if os.Getenv("ENV_PATH") == "" {
		os.Setenv("ENV_PATH", ".env")
	}

	// load environment variables
	if err := barf.Env(primer.ENV, os.Getenv("ENV_PATH")); err != nil {
		barf.Logger().Fatalf(`[main.main] [barf.Env(primer.ENV, os.Getenv("ENV_PATH"))] %s`, err.Error())
	}

	// configure barf
	if err := barf.Stark(barf.Augment{
		Port:         primer.ENV.Port,
		Logging:      barf.Allow(), // enable request logging
		Recovery:     barf.Allow(), // enable panic recovery
		WriteTimeout: 30,
		ReadTimeout:  30,
		CORS: &barf.CORS{
			AllowedOrigins: []string{
				"*",
			},
			AllowedMethods: []string{
				http.MethodGet,
				http.MethodPost,
				http.MethodPut,
				http.MethodPatch,
				http.MethodDelete,
				http.MethodOptions,
			},
		},
	}); err != nil {
		barf.Logger().Fatalf(`[main.main] [barf.Stark(barf.Augment)] %s`, err.Error())
	}

	// connect database
	if err := database.NewPostgreSQLConnection(primer.ENV.PostgreSQLURI, primer.ENV.PostgreSQLConnections, primer.ENV.PostgreSQLDebug); err != nil {
		barf.Logger().Fatalf(`[main.main] [database.NewPostgreSQLConnection(primer.ENV.PostgreSQLURI, primer.ENV.PostgreSQLConnections, primer.ENV.PostgreSQLDebug)] %s`, err.Error())
	}

	if err := migration.CreateTables(); err != nil {
		barf.Logger().Fatalf(`[main.main] [database.CreateTables()] %s`, err.Error())
	}

	// if err := database.ReadFileAndExecuteQueries(primer.ENV.SQLFilePath); err != nil {
	// 	barf.Logger().Fatalf(`[main.main] [database.ReadFileAndExecuteQueries(primer.ENV.SQLFilePath)] %s`, err.Error())
	// }

	// preload v1 routes
	version.V1()

	// call upon barf to listen and serve
	if err := barf.Beck(); err != nil {
		barf.Logger().Errorf(`[main.main] [barf.Beck()] %s`, err.Error())
		os.Exit(1)
	}
}
