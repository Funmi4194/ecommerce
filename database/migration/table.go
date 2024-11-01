package migration

import (
	"context"

	"github.com/funmi4194/ecommerce/database"
	orderRepository "github.com/funmi4194/ecommerce/repository/order"
	productRepository "github.com/funmi4194/ecommerce/repository/product"
	userRepository "github.com/funmi4194/ecommerce/repository/user"
	"github.com/opensaucerer/barf"
)

var Table = []interface{}{
	&userRepository.User{},
	&orderRepository.Order{},
	&productRepository.Product{},
}

// CreateTables creates tables that do not already exist. Although we have connections to other DBs configure.Save should only handle migration for configure.Save DB.
func CreateTables() error {
	for _, m := range Table {
		_, err := database.PostgreSQLDB.NewCreateTable().
			IfNotExists().
			Model(m).Exec(context.TODO())
		if err != nil {
			barf.Logger().Warnf("failed to create %v table", m)
			return err
		}
	}
	return nil
}

// migrate effects any database schema migration
func Migrate() error {
	return nil
}
