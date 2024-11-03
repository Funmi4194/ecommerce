package order

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"strings"
	"time"

	"github.com/funmi4194/ecommerce/database"
	"github.com/funmi4194/ecommerce/enum"
	"github.com/funmi4194/ecommerce/reflection"
	"github.com/funmi4194/ecommerce/types"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/schema"
)

// Scan implements the Scanner interface.
func (i *Item) Scan(src interface{}) error {
	switch v := src.(type) {
	case []byte:
		return json.Unmarshal(v, i)
	case string:
		return json.Unmarshal([]byte(v), i)
	case nil:
		return nil
	}
	return nil
}

// Value implements the driver Valuer interface.
func (i Item) Value() (driver.Value, error) {
	b, err := json.Marshal(i)
	return string(b), err
}

/* Fields returns the struct fields as a slice of interface{} values */
func (o *Order) Fields() []interface{} {
	return reflection.ReturnStructFields(o)
}

/*
Date loads the created_at and updated_at fields of the transaction if not already present, otherwise, it loads the updated_at field only.

If the "pessimistic" parameter is set to true, it loads both fields regardless
*/
func (o *Order) Date(pessimistic ...bool) {
	if len(pessimistic) > 0 && !pessimistic[0] {
		if o.CreatedAt.IsZero() {
			o.CreatedAt = schema.NullTime{Time: time.Now()}
			o.UpdatedAt = schema.NullTime{Time: time.Now()}
			return
		}
		o.UpdatedAt = schema.NullTime{Time: time.Now()}
		return
	}
	o.CreatedAt = schema.NullTime{Time: time.Now()}
	o.UpdatedAt = schema.NullTime{Time: time.Now()}
}

/*
FByMap finds and returns a orders matching the key/value pairs provided in the map

By default, only the id and user_id fields are loaded

The	"preloadandjoin" parameter can be used to request that all the fields of the struct be loaded

It returns an error if any
*/
func (o *Order) FByMap(m types.SQLMaps, preloadandjoin ...bool) error {
	query, args := database.MapsToWQuery(m)
	if len(preloadandjoin) > 1 && preloadandjoin[0] && preloadandjoin[1] {
		return database.PostgreSQLDB.NewRaw(`SELECT * FROM orders WHERE `+query, args...).Scan(context.Background(), o)
	}
	if len(preloadandjoin) > 0 && preloadandjoin[0] {
		return database.PostgreSQLDB.NewRaw(`SELECT * FROM orders WHERE `+query, args...).Scan(context.Background(), o)
	}
	return database.PostgreSQLDB.NewRaw(`SELECT id, user_id FROM orders WHERE `+query, args...).Scan(context.Background(), o)
}

/*
Create inserts a new order into the database using the provided transaction

It returns an error if any
*/
func (o *Order) CreateTx(tx *bun.Tx, m types.SQLMaps) error {
	query, args := database.MapsToIQuery(m)
	if _, err := tx.NewRaw(`INSERT INTO orders `+query, args...).Exec(context.Background()); err != nil && err != sql.ErrNoRows {
		return err
	}
	return nil
}

/*
FUByMap finds and returns a orders matching the key/value pairs provided in the map for the purpose of an update thereby causing the matching rows to be locked

By default, only the id and user_id fields are loaded

The	"preloadandjoin" parameter can be used to request that all the fields of the struct be loaded

It returns an error if any
*/
func (o *Order) FUByMap(tx *bun.Tx, m types.SQLMaps, preloadandjoin ...bool) error {
	query, args := database.MapsToWQuery(m)
	if len(preloadandjoin) > 1 && preloadandjoin[0] && preloadandjoin[1] {
		return tx.NewRaw(`SELECT * FROM orders WHERE `+query+` FOR UPDATE`, args...).Scan(context.Background(), o)
	}
	if len(preloadandjoin) > 0 && preloadandjoin[0] {
		return tx.NewRaw(`SELECT * FROM orders WHERE `+query+` FOR UPDATE`, args...).Scan(context.Background(), o)
	}
	return tx.NewRaw(`SELECT id, user_id FROM orders WHERE `+query+` FOR UPDATE`, args...).Scan(context.Background(), o)
}

/*
UByMapTx updates a orders matching the key/value pairs provided in the map using the provided transaction

It returns an error if any
*/
func (o *Order) UByMapTx(tx *bun.Tx, m types.SQLMaps) error {
	query, args := database.MapsToSQuery(m)
	if strings.Contains(query, string(enum.RETURNING)) {
		return tx.NewRaw(`UPDATE orders `+query, args...).Scan(context.Background(), o)
	}
	_, err := tx.NewRaw(`UPDATE orders `+query, args...).Exec(context.Background())
	return err
}

/*
FByMap finds and returns all orders matching the key/value pairs provided in the map

By default, only the id and user_id fields are loaded

The	"preloadandjoin" parameter can be used to request that all the fields of the struct be loaded

It returns an error if any
*/
func (o *Orders) FByMap(m types.SQLMaps, limit, offset int, sort string, preloadandjoin ...bool) error {
	query, args := database.MapsToWQuery(m)
	if len(preloadandjoin) > 1 && preloadandjoin[0] && preloadandjoin[1] {
		if query != "" {
			query = `SELECT * FROM orders WHERE ` + query + ` ORDER BY orders.updated_at ` + sort
		} else {
			query = `SELECT * FROM orders ORDER BY orders.updated_at ` + sort
		}

		if limit > 0 {
			query += ` LIMIT ?`
		}

		if offset > 0 {
			query += ` OFFSET ?`
		}

		rows, err := database.PostgreSQLDB.QueryContext(context.Background(), query, append(args, limit, offset)...)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var order Order
			if err := rows.Scan(order.Fields()...); err != nil {
				return err
			}
			*o = append(*o, order)
		}
	}
	if len(preloadandjoin) > 0 && preloadandjoin[0] {
		if query != "" {
			query = `SELECT * FROM orders WHERE ` + query + ` ORDER BY orders.updated_at ` + sort
		} else {
			query = `SELECT * FROM orders ORDER BY orders.updated_at ` + sort
		}

		if limit > 0 {
			query += ` LIMIT ?`
		}

		if offset > 0 {
			query += ` OFFSET ?`
		}

		return database.PostgreSQLDB.NewRaw(query, append(args, limit, offset)...).Scan(context.Background(), o)
	}
	if query != "" {
		query = `SELECT id, user_id FROM orders WHERE ` + query + ` ORDER BY orders.updated_at ` + sort
	} else {
		query = `SELECT id, user_id FROM orders ORDER BY orders.updated_at ` + sort
	}

	if limit > 0 {
		query += ` LIMIT ?`
	}

	if offset > 0 {
		query += ` OFFSET ?`
	}

	return database.PostgreSQLDB.NewRaw(query, append(args, limit, offset)...).Scan(context.Background(), o)
}

/*
CByMap finds and counts all orders matching the key/value pairs provided in the map

It returns an error if any
*/
func (o *Orders) CByMap(m types.SQLMaps) (int, error) {
	var count int
	query, args := database.MapsToWQuery(m)
	if query != "" {
		query = `SELECT count(*) FROM orders WHERE ` + query
	} else {
		query = `SELECT count(*) FROM orders`
	}
	err := database.PostgreSQLDB.NewRaw(query, args...).Scan(context.Background(), &count)
	return count, err
}
