package product

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/funmi4194/ecommerce/enum"
	"github.com/funmi4194/ecommerce/helper"
	"github.com/funmi4194/ecommerce/primer"
	commonRepository "github.com/funmi4194/ecommerce/repository/common"
	productRepository "github.com/funmi4194/ecommerce/repository/product"
	userRepository "github.com/funmi4194/ecommerce/repository/user"
	"github.com/funmi4194/ecommerce/types"
	"github.com/opensaucerer/barf"
	"github.com/uptrace/bun"
)

func Publish(userId string, payload types.Publish) (*productRepository.Products, error) {

	// create a new transaction
	btx, err := commonRepository.BeginTx()
	if err != nil {
		return nil, err
	}
	defer btx.Rollback()

	user := userRepository.User{
		ID: userId,
	}

	// find user by Id
	err = user.FByKeyVal("id", user.ID, true)
	if err != nil {
		barf.Logger().Errorf(`[product.Publish] [user.FByKeyVal("id", user.ID, true)] %s`, err.Error())
		if err == sql.ErrNoRows {
			return nil, errors.New("looks like your account no longer exists. please contact support")
		}
		return nil, errors.New("we're having issues publishing products. please try again later")
	}

	if user.Role != enum.Admin {
		return nil, errors.New("you do not have the permission to this feature")
	}

	var product productRepository.Product
	products := make(productRepository.Products, 0)

	insertMap := types.SQLMaps{
		IMaps: []types.SQLMap{},
	}

	for _, p := range payload.Products {

		// ensure name is provided
		if p.Name == "" {
			return nil, errors.New("product name is required")
		}

		// ensure user is not uploading a product with 0 in stock
		if p.Stock <= 0 {
			return nil, errors.New("product stock should be greater than zero")
		}

		// ensure user provided the product image
		if p.ProductUrl == "" {
			return nil, errors.New("product image is required")
		}

		products = append(products, productRepository.Product{
			ID:          helper.GenerateUUID(),
			Name:        p.Name,
			Price:       p.Price,
			Stock:       p.Stock,
			ProductUrl:  p.ProductUrl,
			Status:      enum.Published,
			Description: p.Description,
			CreatedAt:   bun.NullTime{Time: time.Now()},
			UpdatedAt:   bun.NullTime{Time: time.Now()},
		})

		insertMap.IMaps = append(insertMap.IMaps, types.SQLMap{
			Map: map[string]interface{}{
				"id":          helper.GenerateUUID(),
				"name":        p.Name,
				"price":       p.Price,
				"stock":       p.Stock,
				"product_url": p.ProductUrl,
				"status":      enum.Published,
				"description": p.Description,
				"created_at":  bun.NullTime{Time: time.Now()},
				"updated_at":  bun.NullTime{Time: time.Now()},
			},
		})
	}

	if len(insertMap.IMaps) > 0 {
		if err := product.CreateTx(btx, insertMap); err != nil {
			barf.Logger().Errorf(`[product.Publish] [product.CreateTx(btx, insertMap)] %s`, err.Error())
			return nil, errors.New("we're having issues publishing products. please try again later")
		}
	}

	// commit transaction
	if err := btx.Commit(); err != nil {
		barf.Logger().Errorf(`[product.Publish] [btx.Commit()] %s`, err.Error())
		return nil, errors.New("we're having issues publishing products. please try again later")
	}

	return &products, nil
}

func UpdateProduct(userId string, payload types.UpdateProduct) (*productRepository.Product, error) {

	// create a new transaction
	btx, err := commonRepository.BeginTx()
	if err != nil {
		return nil, err
	}
	defer btx.Rollback()

	user := userRepository.User{
		ID: userId,
	}

	// find user by Id
	err = user.FByKeyVal("id", user.ID, true)
	if err != nil {
		barf.Logger().Errorf(`[product.Publish] [user.FByKeyVal("id", user.ID, true)] %s`, err.Error())
		if err == sql.ErrNoRows {
			return nil, errors.New("looks like your account no longer exists. please contact support")
		}
		return nil, errors.New("we're having issues updating product. please try again later")
	}

	if user.Role != enum.Admin {
		return nil, errors.New("you do not have the permission to this feature")
	}

	if payload.ProductId == "" {
		return nil, errors.New("product id is required")
	}

	var product productRepository.Product

	// do an update if the product already exists
	err = product.FByKeyVal("id", payload.ProductId, true)
	if err != nil {
		if err == sql.ErrNoRows {
			barf.Logger().Errorf(`[product.UpdateProduct] [product.FByKeyVal("id", payload.ProductId)] %s`, err.Error())
			return nil, errors.New("product item not found")
		}
		return nil, errors.New("we're having issues updating product. please try again later")
	}

	//  generate filter map
	query := map[string]interface{}{
		"updated_at": bun.NullTime{Time: time.Now()},
	}

	if payload.Status != nil {
		query["status"] = &payload.Status
	}

	if payload.Name != nil {
		query["name"] = &payload.Name
	}

	if payload.Price != nil {
		query["price"] = &payload.Price
	}

	if payload.Stock != nil {
		query["stock"] = &payload.Stock
	}

	if payload.ProductUrl != nil {
		query["product_url"] = &payload.ProductUrl
	}

	if payload.Description != nil {
		query["description"] = &payload.Description
	}

	if err := product.UByMapTx(btx, types.SQLMaps{
		WMaps: []types.SQLMap{
			{
				Map: map[string]interface{}{
					"products.id": payload.ProductId,
				},
				JoinOperator:       enum.And,
				ComparisonOperator: enum.Equal,
			},
		},
		SMap: types.SQLMap{
			Map:                query,
			JoinOperator:       enum.Comma,
			ComparisonOperator: enum.Equal,
		},
		RMap: types.SQLMap{
			Map: map[string]interface{}{"*": nil},
		},
		WJoinOperator: enum.And,
	}); err != nil {
		barf.Logger().Errorf(`[product.UpdateProduct] [product.UByMapTx(btx, types.SQLMaps{s] %s`, err.Error())
		return nil, errors.New("we're having issues updating product. please try again later")
	}

	// commit transaction
	if err := btx.Commit(); err != nil {
		barf.Logger().Errorf(`[user.Deactivate] [btx.Commit()] %s`, err.Error())
		return nil, errors.New("we're having issues updating product. please try again later")
	}

	return &product, nil
}

// Product retrieve a single product by Id
func Product(userId string, payload types.ProductFilter) (*productRepository.Product, error) {

	user := userRepository.User{
		ID: userId,
	}

	// find user by Id
	err := user.FByKeyVal("id", user.ID, true)
	if err != nil {
		barf.Logger().Errorf(`[product.Product] [user.FByKeyVal("id", user.ID, true)] %s`, err.Error())
		if err == sql.ErrNoRows {
			return nil, errors.New("looks like your account no longer exists. please contact support")
		}
		return nil, errors.New("we're having issues getting product. please try again later")
	}

	if user.Role != enum.Admin {
		return nil, errors.New("you do not have the permission to this feature")
	}

	if payload.ProductId == "" {
		return nil, errors.New("product id is required")
	}

	var product productRepository.Product

	// do an update if the product already exists
	err = product.FByKeyVal("id", payload.ProductId, true)
	if err != nil {
		if err == sql.ErrNoRows {
			barf.Logger().Errorf(`[product.Product] [product.FByKeyVal("id", payload.ProductId, true)] %s`, err.Error())
			return nil, errors.New("product item not found")
		}
		return nil, errors.New("we're having issues getting product. please try again later")
	}

	return &product, nil
}

// Products retrieve all products
func Products(userId string, payload types.ProductFilter) (*productRepository.Products, *commonRepository.Pagination, error) {

	user := userRepository.User{
		ID: userId,
	}

	// find user by Id
	err := user.FByKeyVal("id", user.ID, true)
	if err != nil {
		barf.Logger().Errorf(`[product.Products] [user.FByKeyVal("id", user.ID, true)] %s`, err.Error())
		if err == sql.ErrNoRows {
			return nil, nil, errors.New("looks like your account no longer exists. please contact support")
		}
		return nil, nil, errors.New("we're having issues updating product. please try again later")
	}

	//  generate filter map
	EqFilter := map[string]interface{}{}
	gtEqFilter := map[string]interface{}{}
	ltEqFilter := map[string]interface{}{}
	searchFilter := map[string]interface{}{}

	// if user is an admin support status filter
	if user.Role == enum.Admin {
		if payload.Status != "" {
			fmt.Println("hi")
			EqFilter["status"] = payload.Status
		}
	} else {
		EqFilter["status"] = enum.Published
	}

	if payload.ProductId != "" {
		EqFilter["id"] = payload.ProductId
	}

	if payload.MinAmount != nil {
		gtEqFilter["CAST(price AS NUMERIC)"] = enum.SQLRaw{
			Value: fmt.Sprintf("CAST(price AS NUMERIC) >= %d", *payload.MinAmount),
		}
	}

	if payload.MaxAmount != nil {
		ltEqFilter["CAST(price AS NUMERIC)"] = enum.SQLRaw{
			Value: fmt.Sprintf("CAST(price AS NUMERIC) <= %d", *payload.MaxAmount),
		}
	}

	if !payload.StartDate.IsZero() {
		gtEqFilter["created_at"] = payload.StartDate
	}
	if !payload.EndDate.IsZero() {
		ltEqFilter["created_at"] = payload.EndDate
	}

	if payload.Search != "" {
		searchFilter["name"] = fmt.Sprintf("%%%s%%", payload.Search)
		searchFilter["description"] = fmt.Sprintf("%%%s%%", payload.Search)
	}

	orderMap := map[string]interface{}{
		"products.updated_at": "DESC",
	}

	limit := primer.PageLimit
	page := 1
	offset := 0

	if payload.Limit != nil && *payload.Limit > 0 {
		limit = *payload.Limit
	}

	if payload.Page != nil && *payload.Page > 0 {
		offset = (*payload.Page - 1) * limit
		page = *payload.Page
	}

	queryMap := []types.SQLMap{
		{
			Map:                EqFilter,
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
		{
			Map:                searchFilter,
			JoinOperator:       enum.Or,
			ComparisonOperator: enum.ILike,
		},
	}

	products := make(productRepository.Products, 0)

	if err := products.FByMap(types.SQLMaps{
		WMaps: queryMap,
		OMap: types.SQLMap{
			Map: orderMap,
		},
		WJoinOperator: enum.And,
	}, limit, offset, enum.DESC.String(), true, true); err != nil {
		barf.Logger().Errorf(`[product.Products] [products.FByMap(types.SQLMaps{] %s`, err.Error())
		if err == sql.ErrNoRows {
			return nil, nil, errors.New("we couldn't find any products")
		}
		return nil, nil, errors.New("we're having issues retrieving products. please try again later")
	}

	var pagination = &commonRepository.Pagination{}

	if payload.Paginate {
		total, err := products.CByMap(types.SQLMaps{
			WMaps:         queryMap,
			WJoinOperator: enum.And,
		})
		if err != nil {
			barf.Logger().Errorf(`[product.Products] [product.Products].CByMap(types.SQLMaps{] %s`, err.Error())
			if err == sql.ErrNoRows {
				return nil, nil, errors.New("products not found")
			}
			return nil, nil, errors.New("we're having issues retrieving products. please try again later")
		}

		pagination = &commonRepository.Pagination{
			Page:  page,
			Limit: limit,
			Total: total,
			Pages: int(math.Ceil(float64(total) / float64(limit))),
		}
	}

	return &products, pagination, nil
}
