package types

import "github.com/funmi4194/ecommerce/primitive"

type Env struct {
	// Port for the server to listen on
	Port string `barfenv:"key=PORT;required=true"`
	// Number of connections to the database
	PostgreSQLConnections int32 `barfenv:"key=POSTGRESQL_CONNECTIONS;required=true"`
	// Database connection string
	PostgreSQLURI string `barfenv:"key=POSTGRESQL_URI;required=true"`
	// Path to SQL file containing queries to be executed on startup
	SQLFilePath string `barfenv:"key=SQL_FILE_PATH;required=true"`
	// Enables verbose logging of database queries
	PostgreSQLDebug bool `barfenv:"key=POSTGRESQL_DEBUG;required=true"`
	// Name of the app instance
	AppName primitive.String `barfenv:"key=APP_NAME;required=true"`
	// Secret for generating JWT signatures
	JWTSecret string `barfenv:"key=JWT_SECRET;required=true"`
	// GoogleApplicationCredentials is the path to the google application credentials
	GoogleApplicationCredentials string `barfenv:"key=GOOGLE_APPLICATION_CREDENTIALS;required=true"`
}
