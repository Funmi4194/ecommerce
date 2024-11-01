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
	//Email is the default email address for the app
	EmailAddress string `barfenv:"key=EMAIL_ADDRESS;required=true"`
	//App password for the default email address
	EmailPassword string `barfenv:"key=EMAIL_PASSWORD;required=true"`
	// Email Port is Port for stmp host
	EmailPort string `barfenv:"key=EMAIL_PORT;required=true"`
	// Email SMTP is the smtp host
	EmailSMTP string `barfenv:"key=EMAIL_SMTP;required=true"`
}
