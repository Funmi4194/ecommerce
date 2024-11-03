package order

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/funmi4194/ecommerce/enum"
	"github.com/funmi4194/ecommerce/helper"
	"github.com/funmi4194/ecommerce/primer"
	"github.com/funmi4194/ecommerce/primitive"
	commonRepository "github.com/funmi4194/ecommerce/repository/common"
	orderRepository "github.com/funmi4194/ecommerce/repository/order"
	productRepository "github.com/funmi4194/ecommerce/repository/product"
	userRepository "github.com/funmi4194/ecommerce/repository/user"
	"github.com/funmi4194/ecommerce/types"
	"github.com/opensaucerer/barf"
	"github.com/uptrace/bun/schema"
)

// InitiateOrder is the logic function to create an order for a user
func InitiateOrder(userId string, payload types.InitiateOrder) (*orderRepository.Order, error) {

	if len(payload.Items) == 0 {
		return nil, errors.New("you need to select at least one product to order")
	}

	if len(payload.Items) > 100 {
		return nil, errors.New("you can only purchase up to 100 products at a time")
	}

	// create a new transaction
	btx, err := commonRepository.BeginTx()
	if err != nil {
		return nil, err
	}
	defer btx.Rollback()

	user := userRepository.User{
		ID: userId,
	}

	// find user by ID
	err = user.FByKeyVal("id", user.ID, true)
	if err != nil {
		barf.Logger().Errorf(`[order.InitiateOrder] [user.FByKeyVal("id", user.ID, true)] %s`, err.Error())
		if err == sql.ErrNoRows {
			return nil, errors.New("looks like your account no longer exists. please contact support")
		}
		return nil, errors.New("we're having issues initiating order. please try again later")
	}

	// verify the items exist
	itemIds := []interface{}{}
	for _, item := range payload.Items {
		if item.ProductId != "" {
			itemIds = append(itemIds, item.ProductId)
		}
	}

	products := make(productRepository.Products, 0)
	o := orderRepository.Order{}

	// get all products
	if err := products.FByMap(types.SQLMaps{
		WMaps: []types.SQLMap{
			{
				Map: map[string]interface{}{
					"id": itemIds,
				},
				JoinOperator:       enum.And,
				ComparisonOperator: enum.In,
			},
			{
				Map: map[string]interface{}{
					"status": enum.Published,
				},
				JoinOperator:       enum.And,
				ComparisonOperator: enum.Equal,
			},
		},
		WJoinOperator: enum.And,
	}, 100, 0, enum.DESC.String(), true, true); err != nil {
		barf.Logger().Errorf(`[creation.Checkout] [items.FByMap(types.SQLMaps{] %s`, err.Error())
		if err == sql.ErrNoRows {
			return nil, errors.New("we couldn't find any of the products you're trying to checkout")
		}
		return nil, errors.New("we're having issues checking out those products. please try again later")
	}

	if len(products) != len(itemIds) {
		return nil, errors.New("some of the selected products are no longer available. please refresh and try again")
	}

	if len(products) == 1 {
		o.ProductID = products[0].ID
	}

	for _, product := range products {

		for _, item := range payload.Items {
			if item.ProductId == product.ID {

				if item.Quantity <= 0 {
					return nil, errors.New("item quantity is required")
				}

				if item.Quantity > int(product.Stock) {
					return nil, fmt.Errorf("requested quantity for product '%s' exceeds available stock", product.Name)
				}

				o.Invoice = append(o.Invoice, orderRepository.Item{
					Key:      product.ID,
					Name:     product.Name,
					Amount:   product.Price,
					Quantity: item.Quantity,
				})
			}
		}
	}

	// compute checksum - this helps prevent duplicate invoice for a tx without paramter changes
	o.Checksum = primer.StringSha256(primer.Stringify(o.Invoice))

	order := orderRepository.Order{}

	// find order
	err = order.FByMap(types.SQLMaps{
		WMaps: []types.SQLMap{
			{
				Map: map[string]interface{}{
					"paid":      false,
					"failed":    false,
					"cancelled": false,
					"checksum":  o.Checksum,
				},
				JoinOperator:       enum.And,
				ComparisonOperator: enum.Equal,
			},
		},
		WJoinOperator: enum.And,
	}, true)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	} else if err == nil {
		// if tx exists, return it
		if order.ID != "" {
			return &order, nil
		}
	}

	// compute total amount
	o.Amount = 0
	for _, item := range o.Invoice {
		o.Amount += math.Ceil((item.Amount*float64(item.Quantity))*100) / 100
	}

	o.UserID = user.ID
	o.ID = helper.GenerateUUID()
	o.Status = enum.Pending
	o.Reference = helper.GenerateRef()
	o.Date()
	o.History = []commonRepository.History{
		{
			Act: "Initiated product(s) purchase",
			By:  user.ID,
			At:  o.CreatedAt,
		},
	}
	o.Remark = "Product(s) purchase"

	// create order
	if err := o.CreateTx(btx, types.SQLMaps{
		IMaps: []types.SQLMap{
			{
				Map: map[string]interface{}{
					"id":           o.ID,
					"user_id":      o.UserID,
					"status":       o.Status,
					"reference":    o.Reference,
					"paid":         o.Paid,
					"paid_at":      o.PaidAt,
					"cancelled":    o.Cancelled,
					"cancelled_at": o.CancelledAt,
					"failed":       o.Failed,
					"failed_at":    o.FailedAt,
					"checksum":     o.Checksum,
					"history":      o.History,
					"invoice":      o.Invoice,
					"amount":       o.Amount,
					"remark":       o.Remark,
					"product_id":   o.ProductID,
					"created_at":   o.CreatedAt,
					"updated_at":   o.UpdatedAt,
				},
			},
		},
	}); err != nil {
		return nil, err
	}

	// commit transaction
	if err := btx.Commit(); err != nil {
		barf.Logger().Errorf(`[creation.Checkout] [btx.Commit()] %s`, err.Error())
		return nil, errors.New("we're having issues checking out those assets. please try again later")
	}

	return &o, nil
}

// CancelOrder is the logic function to cancel order
func CancelOrder(userId string, payload types.CancelOrder) (*orderRepository.Order, error) {

	// start a database transaction
	btx, err := commonRepository.BeginTx()
	if err != nil {
		return nil, err
	}
	defer btx.Rollback()

	user := userRepository.User{
		ID: userId,
	}

	// find user by ID
	err = user.FByKeyVal("id", user.ID, true)
	if err != nil {
		barf.Logger().Errorf(`[order.CancelOrder] [user.FByKeyVal("id", user.ID, true)] %s`, err.Error())
		if err == sql.ErrNoRows {
			return nil, errors.New("looks like your account no longer exists. please contact support")
		}
		return nil, errors.New("we're having issues cancelling the order. please try again later")
	}

	var order orderRepository.Order

	if payload.OrderId == "" {
		return nil, errors.New("order id is required")
	}

	// find order and lock
	err = order.FUByMap(btx, types.SQLMaps{
		WMaps: []types.SQLMap{
			{
				Map: map[string]interface{}{
					"id":      payload.OrderId,
					"user_id": user.ID,
				},
				JoinOperator:       enum.And,
				ComparisonOperator: enum.Equal,
			},
		},
		WJoinOperator: enum.And,
	}, true)
	if err != nil {
		barf.Logger().Errorf(`[order.CancelOrder] [order.FUByMap(btx, types.SQLMaps{] %s`, err.Error())
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("order not found")
		}
		return nil, errors.New("we're having issues cancelling the order. please try again later")
	}

	if payload.Cancel {

		if order.Cancelled {
			return nil, errors.New("order has already been cancelled")
		}

		// only allow user cancel when order is still pending
		if order.Status != enum.Pending {
			return nil, errors.New("order cannot be cancelled")
		}

		// update order
		err = order.UByMapTx(btx, types.SQLMaps{
			WMaps: []types.SQLMap{
				{
					Map: map[string]interface{}{
						"id":      order.ID,
						"user_id": user.ID,
					},
					JoinOperator:       enum.And,
					ComparisonOperator: enum.Equal,
				},
			},

			SMap: types.SQLMap{
				Map: map[string]interface{}{
					"cancelled":    true,
					"cancelled_at": "now()",
					"remark":       "Order cancelled",
					"status":       enum.Cancelled.String(),
					"history": enum.SQLValueMerge{
						Operator: enum.CONCAT,
						Values: primitive.Array{
							commonRepository.History{
								Act: "Cancelled order",
								By:  user.ID,
								At:  schema.NullTime{Time: time.Now()},
							},
						},
					},
					"updated_at": "now()", // update updated_at
				},
				JoinOperator:       enum.Comma,
				ComparisonOperator: enum.Equal,
			},

			RMap: types.SQLMap{
				Map: map[string]interface{}{"*": nil},
			},

			WJoinOperator: enum.And,
		})
		if err != nil {
			barf.Logger().Errorf(`[order.CancelOrder] [order.UByMapTx(btx, types.SQLMaps{] %s`, err.Error())
			if err == sql.ErrNoRows {
				return nil, errors.New("order not found")
			}
			return nil, err
		}
	}

	// commit transaction
	if err := btx.Commit(); err != nil {
		barf.Logger().Errorf(`[order.CancelOrder] [btx.Commit()] %s`, err.Error())
		return nil, errors.New("we're having issues cancelling order. please try again later")
	}

	return &order, nil
}

func UpdateOrder(userId string, payload types.UpdateOrder) (*orderRepository.Order, error) {

	// start a database transaction
	btx, err := commonRepository.BeginTx()
	if err != nil {
		return nil, err
	}
	defer btx.Rollback()

	user := userRepository.User{
		ID: userId,
	}

	// find user by ID
	err = user.FByKeyVal("id", user.ID, true)
	if err != nil {
		barf.Logger().Errorf(`[order.UpdateOrder] [user.FByKeyVal("id", user.ID, true)] %s`, err.Error())
		if err == sql.ErrNoRows {
			return nil, errors.New("looks like your account no longer exists. please contact support")
		}
		return nil, errors.New("we're having issues updating the order. please try again later")
	}

	if user.Role != enum.Admin {
		return nil, errors.New("you do not have the permission to access this feature")
	}

	var order orderRepository.Order

	if payload.OrderId == "" {
		return nil, errors.New("order id is required")
	}

	// find order and lock
	err = order.FUByMap(btx, types.SQLMaps{
		WMaps: []types.SQLMap{
			{
				Map: map[string]interface{}{
					"id": payload.OrderId,
				},
				JoinOperator:       enum.And,
				ComparisonOperator: enum.Equal,
			},
		},
		WJoinOperator: enum.And,
	}, true)
	if err != nil {
		barf.Logger().Errorf(`[order.UpdateOrder] [order.FUByMap(btx, types.SQLMaps{] %s`, err.Error())
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("order not found")
		}
		return nil, errors.New("we're having issues updating the order. please try again later")
	}

	// if order has been cnacelled be a user admin should not be able to update order
	if order.Cancelled {
		return nil, errors.New("order has already been cancelled")
	}

	if payload.Status != enum.Completed && payload.Status != enum.Approved && payload.Status != enum.Rejected && payload.Status != enum.Cancelled {
		return nil, errors.New("order can either be completed, approved, rejected or cancelled")
	}

	// update order
	err = order.UByMapTx(btx, types.SQLMaps{
		WMaps: []types.SQLMap{
			{
				Map: map[string]interface{}{
					"id": order.ID,
				},
				JoinOperator:       enum.And,
				ComparisonOperator: enum.Equal,
			},
		},

		SMap: types.SQLMap{
			Map: map[string]interface{}{
				"status":     payload.Status,
				"updated_at": "now()", // update updated_at
			},
			JoinOperator:       enum.Comma,
			ComparisonOperator: enum.Equal,
		},

		RMap: types.SQLMap{
			Map: map[string]interface{}{"*": nil},
		},

		WJoinOperator: enum.And,
	})
	if err != nil {
		barf.Logger().Errorf(`[order.UpdateOrder] [order.UByMapTx(btx, types.SQLMaps{] %s`, err.Error())
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	// commit transaction
	if err := btx.Commit(); err != nil {
		barf.Logger().Errorf(`[order.CancelOrder] [btx.Commit()] %s`, err.Error())
		return nil, errors.New("we're having issues cancelling order. please try again later")
	}

	return &order, nil
}

func Orders(userId string, payload types.OrderFilter) (*orderRepository.Orders, *commonRepository.Pagination, error) {
	user := userRepository.User{
		ID: userId,
	}

	// find user by ID
	err := user.FByKeyVal("id", user.ID, true)
	if err != nil {
		barf.Logger().Errorf(`[order.Orders] [user.FByKeyVal("id", user.ID, true)] %s`, err.Error())
		if err == sql.ErrNoRows {
			return nil, nil, errors.New("looks like your account no longer exists. please contact support")
		}
		return nil, nil, errors.New("we're having issues retrieving orders. please try again later")
	}

	orders := make(orderRepository.Orders, 0)

	//  generate filter map
	Eqfilter := map[string]interface{}{}

	// verify user is an admin
	if user.Role != enum.Admin {
		Eqfilter["user_id"] = user.ID
	}

	if payload.Paid != nil {
		Eqfilter["paid"] = *payload.Paid
	}
	if payload.Cancelled != nil {
		Eqfilter["cancelled"] = *payload.Cancelled
	}
	if payload.Failed != nil {
		Eqfilter["failed"] = *payload.Failed
	}
	if payload.OrderId != "" {
		Eqfilter["id"] = payload.OrderId
	}
	if payload.Reference != "" {
		Eqfilter["reference"] = payload.Reference
	}

	gtEqFilter := map[string]interface{}{}
	ltEqFilter := map[string]interface{}{}

	if payload.MinAmount != nil {
		gtEqFilter["amount"] = *payload.MinAmount
	}
	if payload.MaxAmount != nil {
		ltEqFilter["amount"] = *payload.MaxAmount
	}
	if !payload.StartDate.IsZero() {
		gtEqFilter["created_at"] = payload.StartDate
	}
	if !payload.EndDate.IsZero() {
		ltEqFilter["created_at"] = payload.EndDate
	}

	limit := primer.PageLimit
	page := 1
	offset := 0

	if payload.Limit != nil {
		limit = *payload.Limit
	}

	if payload.Page != nil {
		offset = (*payload.Page - 1) * limit
		page = *payload.Page
	}

	queryMap := []types.SQLMap{
		{
			Map:                Eqfilter,
			JoinOperator:       enum.And,
			ComparisonOperator: enum.Equal,
		},
		{
			Map:                gtEqFilter,
			JoinOperator:       enum.And,
			ComparisonOperator: enum.GreaterThanOrEqual,
		},
		{
			Map:                ltEqFilter,
			JoinOperator:       enum.And,
			ComparisonOperator: enum.LessThanOrEqual,
		},
	}

	err = orders.FByMap(types.SQLMaps{
		WMaps:         queryMap,
		WJoinOperator: enum.And,
	}, limit, offset, enum.DESC.String(), true)
	if err != nil {
		barf.Logger().Errorf(`[order.Orders] [orders.FByMap(types.SQLMaps{] %s`, err.Error())
		if err == sql.ErrNoRows {
			return nil, nil, errors.New("orders not found")
		}
		return nil, nil, errors.New("we're having issues retrieving orders. please try again later")
	}

	var pagination *commonRepository.Pagination

	if payload.Paginate {
		total, err := orders.CByMap(types.SQLMaps{
			WMaps:         queryMap,
			WJoinOperator: enum.And,
		})
		if err != nil {
			barf.Logger().Errorf(`[order.Orders] [orders.CByMap(types.SQLMaps{] %s`, err.Error())
			if err == sql.ErrNoRows {
				return nil, nil, errors.New("orders not found")
			}
			return nil, nil, errors.New("we're having issues retrieving orders. please try again later")
		}

		pagination = &commonRepository.Pagination{
			Page:  page,
			Limit: limit,
			Total: total,
			Pages: int(math.Ceil(float64(total) / float64(limit))),
		}
	}

	return &orders, pagination, nil
}
