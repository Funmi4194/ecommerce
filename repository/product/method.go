package product

import (
	"context"
	"database/sql"
	"strings"

	"github.com/funmi4194/ecommerce/database"
	"github.com/funmi4194/ecommerce/enum"
	"github.com/funmi4194/ecommerce/reflection"
	"github.com/funmi4194/ecommerce/types"
	"github.com/uptrace/bun"
)

/* Fields returns the struct fields as a slice of interface{} values */
func (p *Product) Fields() []interface{} {
	return reflection.ReturnStructFields(p)
}

/*
Create inserts a new product into the database

It returns an error if any
*/
func (p *Product) Create(m types.SQLMaps) error {
	query, args := database.MapsToIQuery(m)
	if _, err := database.PostgreSQLDB.NewRaw(`INSERT INTO products `+query, args...).Exec(context.Background()); err != nil && err != sql.ErrNoRows {
		return err
	}
	return nil
}

/*
Create inserts a new product or products into the database using the provided product

It returns an error if any
*/
func (p *Product) CreateTx(tx *bun.Tx, m types.SQLMaps) error {
	query, args := database.MapsToIQuery(m)
	if _, err := tx.NewRaw(`INSERT INTO products `+query, args...).Exec(context.Background()); err != nil && err != sql.ErrNoRows {
		return err
	}
	return nil
}

/*
FByKeyVal finds and returns a products matching the key/value pair

# By default, only the id and name fields are loaded

The	"preloadandjoin" parameter can be used to request that all the fields of the struct be loaded

It returns an error if any
*/
func (p *Product) FByKeyVal(key string, val interface{}, preloadandjoin ...bool) error {
	if len(preloadandjoin) > 1 && preloadandjoin[0] && preloadandjoin[1] {
		return database.PostgreSQLDB.NewRaw(`SELECT * FROM users WHERE users.`+key+` = ?`, val).Scan(context.Background(), p.Fields()...)
	}
	if len(preloadandjoin) > 0 && preloadandjoin[0] {
		return database.PostgreSQLDB.NewRaw(`SELECT * FROM products WHERE `+key+` = ?`, val).Scan(context.Background(), p)
	}
	return database.PostgreSQLDB.NewRaw(`SELECT id, name FROM products WHERE products.`+key+` = ?`, val).Scan(context.Background(), p)
}

/*
FByKeyVal finds and returns all products matching the key/value pair

# By default, only the id and name fields are loaded

The	"preloadandjoin" parameter can be used to request that all the fields of the struct be loaded

It returns an error if any
*/
func (p *Products) FByKeyVal(key string, val interface{}, limit, offset int, sort string, preloadandjoin ...bool) error {
	if len(preloadandjoin) > 1 && preloadandjoin[0] && preloadandjoin[1] {
		return database.PostgreSQLDB.NewRaw(`SELECT * FROM products WHERE products.`+key+` = ? ORDER BY products.updated_at `+sort+` LIMIT ? OFFSET ?`, val, limit, offset).Scan(context.Background(), preloadandjoin)
	}
	if len(preloadandjoin) > 0 && preloadandjoin[0] {
		return database.PostgreSQLDB.NewRaw(`SELECT * FROM products WHERE `+key+` = ? ORDER BY products.updated_at `+sort+` LIMIT ? OFFSET ?`, val, limit, offset).Scan(context.Background(), p)
	}
	return database.PostgreSQLDB.NewRaw(`SELECT id, name FROM products WHERE `+key+` = ? ORDER BY products.updated_at `+sort+` LIMIT ? OFFSET ?`, val, limit, offset).Scan(context.Background(), p)
}

/*
FUByKeyVal finds and returns a products matching the key/value pair for the purpose of an update thereby causing the matching rows to be locked

# By default, only the id and name fields are loaded

The	"preloadandjoin" parameter can be used to request that all the fields of the struct be loaded

It returns an error if any
*/
func (p *Product) FUByKeyVal(tx *bun.Tx, key string, val interface{}, preloadandjoin ...bool) error {
	if len(preloadandjoin) > 1 && preloadandjoin[0] && preloadandjoin[1] {
		return tx.NewRaw(`SELECT * FROM products WHERE `+key+` = ? FOR UPDATE`, val).Scan(context.Background(), p)
	}
	if len(preloadandjoin) > 0 && preloadandjoin[0] {
		return tx.NewRaw(`SELECT * FROM products WHERE `+key+` = ? FOR UPDATE`, val).Scan(context.Background(), p)
	}
	return tx.NewRaw(`SELECT id, name FROM products WHERE id = ? FOR UPDATE`, p.ID).Scan(context.Background(), p)
}

/*
FUByMap finds and returns a products matching the key/value pairs provided in the map for the purpose of an update thereby causing the matching rows to be locked

# By default, only the id and name fields are loaded

The	"preloadandjoin" parameter can be used to request that all the fields of the struct be loaded

It returns an error if any
*/
func (p *Product) FUByMap(tx *bun.Tx, m types.SQLMaps, preloadandjoin ...bool) error {
	query, args := database.MapsToWQuery(m)
	if len(preloadandjoin) > 1 && preloadandjoin[0] && preloadandjoin[1] {
		return tx.NewRaw(`SELECT * FROM products WHERE `+query+` FOR UPDATE`, args...).Scan(context.Background(), p)
	}
	if len(preloadandjoin) > 0 && preloadandjoin[0] {
		return tx.NewRaw(`SELECT * FROM products WHERE `+query+` FOR UPDATE`, args...).Scan(context.Background(), p)
	}
	return tx.NewRaw(`SELECT id, name FROM products WHERE `+query+` FOR UPDATE`, args...).Scan(context.Background(), p)
}

/*
UByMap updates a products matching the key/value pairs provided in the map

It returns an error if any
*/
func (p *Product) UByMap(m types.SQLMaps) error {
	query, args := database.MapsToSQuery(m)
	if strings.Contains(query, string(enum.RETURNING)) {
		return database.PostgreSQLDB.NewRaw(`UPDATE products `+query, args...).Scan(context.Background(), p)
	}
	_, err := database.PostgreSQLDB.NewRaw(`UPDATE products `+query, args...).Exec(context.Background())
	return err
}

/*
UByMapTx updates a products matching the key/value pairs provided in the map using the provided transaction

It returns an error if any
*/
func (p *Product) UByMapTx(tx *bun.Tx, m types.SQLMaps) error {
	query, args := database.MapsToSQuery(m)
	if strings.Contains(query, string(enum.RETURNING)) {
		return tx.NewRaw(`UPDATE products `+query, args...).Scan(context.Background(), p)
	}
	_, err := tx.NewRaw(`UPDATE products `+query, args...).Exec(context.Background())
	return err
}

/*
CByMap finds and counts all products matching the key/value pairs provided in the map

It returns an error if any
*/
func (p *Product) CByMap(m types.SQLMaps) (int, error) {
	var count int
	query, args := database.MapsToWQuery(m)
	err := database.PostgreSQLDB.NewRaw(`SELECT count(*) FROM products WHERE `+query, args...).Scan(context.Background(), &count)
	return count, err
}

/*
CByMap finds and counts all products matching the key/value pairs provided in the map

It returns an error if any
*/
func (p *Products) CByMap(m types.SQLMaps) (int, error) {
	var count int
	query, args := database.MapsToWQuery(m)
	if query != "" {
		query = `SELECT count(*) FROM products WHERE ` + query
	} else {
		query = `SELECT count(*) FROM products`
	}
	err := database.PostgreSQLDB.NewRaw(query, args...).Scan(context.Background(), &count)
	return count, err
}

/*
DByMap deletes a collection products matching the key/value pairs provided in the map

It returns an error if any
*/
func (p *Product) DByMap(m types.SQLMaps) error {
	query, args := database.MapsToWQuery(m)
	_, err := database.PostgreSQLDB.NewRaw(`DELETE FROM products WHERE `+query, args...).Exec(context.Background())
	return err
}

/*
DByMap deletes a collection products matching the key/value pairs provided in the map

It returns an error if any
*/
func (p *Products) DByMap(m types.SQLMaps) error {
	query, args := database.MapsToWQuery(m)
	_, err := database.PostgreSQLDB.NewRaw(`DELETE FROM products WHERE `+query, args...).Exec(context.Background())
	return err
}

/*
FByMap finds and returns all products matching the key/value pairs provided in the map

# By default, only the id and name fields are loaded

The	"preloadandjoin" parameter can be used to request that all the fields of the struct be loaded

It returns an error if any
*/
func (p *Products) FByMap(m types.SQLMaps, limit, offset int, sort string, preloadandjoin ...bool) error {
	query, args := database.MapsToWQuery(m)
	oquery := database.MapsToOQuery(m)
	if len(m.Args) > 0 {
		args = append(m.Args, args...)
	}
	if len(preloadandjoin) > 1 && preloadandjoin[0] && preloadandjoin[1] {
		if query != "" {
			// join clause can come in here
			query = `SELECT * FROM products WHERE ` + query + oquery + ` LIMIT ? OFFSET ?`
		} else {
			// join clause can come in here
			query = `SELECT * FROM products WHERE ` + oquery + ` LIMIT ? OFFSET ?`
		}
		rows, err := database.PostgreSQLDB.QueryContext(context.Background(), query, append(args, limit, offset)...)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var product Product
			if err := rows.Scan(product.Fields()...); err != nil {
				return err
			}
			*p = append(*p, product)
		}
		return nil
	}
	if len(preloadandjoin) > 0 && preloadandjoin[0] {
		if query != "" {
			query = `SELECT * FROM products WHERE ` + query + ` ORDER BY products.updated_at ` + sort + ` LIMIT ? OFFSET ?`
		} else {
			query = `SELECT * FROM products ORDER BY products.updated_at ` + sort + ` LIMIT ? OFFSET ?`
		}
		return database.PostgreSQLDB.NewRaw(query, append(args, limit, offset)...).Scan(context.Background(), p)
	}
	if query != "" {
		query = `SELECT id, name FROM products WHERE ` + query + ` ORDER BY products.updated_at ` + sort + ` LIMIT ? OFFSET ?`
	} else {
		query = `SELECT id, name FROM products ORDER BY products.updated_at ` + sort + ` LIMIT ? OFFSET ?`
	}
	return database.PostgreSQLDB.NewRaw(query, append(args, limit, offset)...).Scan(context.Background(), p)
}
