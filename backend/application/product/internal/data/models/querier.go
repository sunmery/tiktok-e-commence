// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package models

import (
	"context"
)

type Querier interface {
	// 创建分类数据
	//
	//  INSERT INTO products.categories (name, parent_id)
	//  VALUES ($1, $2)
	//  RETURNING id, name, parent_id, is_active, created_at
	CreateCategories(ctx context.Context, arg CreateCategoriesParams) (ProductsCategories, error)
	// 创建商品
	//
	//  INSERT INTO products.products(name,
	//                                description,
	//                                picture,
	//                                price,
	//                                total_stock)
	//  VALUES ($1, $2, $3, $4, $5)
	//  RETURNING id, name, description, picture, price, total_stock, available_stock, reserved_stock, low_stock_threshold, allow_negative, created_at, updated_at, version
	CreateProduct(ctx context.Context, arg CreateProductParams) (ProductsProducts, error)
	// 关联商品与分类
	// 将商品1关联到分类2（Smartphones）
	//
	//  INSERT INTO products.product_categories (product_id, category_id)
	//  VALUES ($1, $2)
	//  RETURNING product_id, category_id
	CreateProductCategories(ctx context.Context, arg CreateProductCategoriesParams) (ProductsProductCategories, error)
	// 记录库存变更
	//
	//  INSERT INTO products.inventory_history (product_id,
	//                                          old_stock,
	//                                          new_stock,
	//                                          change_reason)
	//  VALUES ($1,
	//          (SELECT total_stock FROM products.products WHERE id = $1),
	//          (SELECT total_stock FROM products.products WHERE id = $1) - $2,
	//          'ORDER_RESERVED')
	//  RETURNING id, product_id, old_stock, new_stock, change_reason, created_at
	CreateProductInventoryHistory(ctx context.Context, arg CreateProductInventoryHistoryParams) (ProductsInventoryHistory, error)
	//GetProduct
	//
	//  SELECT id, name, description, picture, price, total_stock, available_stock, reserved_stock, low_stock_threshold, allow_negative, created_at, updated_at, version
	//  FROM products.products
	//  WHERE id = $1
	//  LIMIT 1
	GetProduct(ctx context.Context, id int32) (ProductsProducts, error)
	// 查询某分类下的所有商品
	//
	//  SELECT p.id, p.name, p.description, p.picture, p.price, p.total_stock, p.available_stock, p.reserved_stock, p.low_stock_threshold, p.allow_negative, p.created_at, p.updated_at, p.version
	//  FROM products.products p
	//           JOIN products.product_categories pc ON p.id = pc.product_id
	//  WHERE pc.category_id = $1
	GetProductCategories(ctx context.Context, categoryID int32) ([]ProductsProducts, error)
	// -- name: CreateAuditLog :one
	// INSERT INTO products.audit_log (action, product_id, owner, name)
	// VALUES ($1, $2, $3, $4)
	// RETURNING *;
	//
	//
	//  SELECT id, name, description, picture, price, total_stock, available_stock, reserved_stock, low_stock_threshold, allow_negative, created_at, updated_at, version
	//  FROM products.products
	//  ORDER BY id
	//  OFFSET $1 LIMIT $2
	ListProducts(ctx context.Context, arg ListProductsParams) ([]ProductsProducts, error)
	//SearchProducts
	//
	//  SELECT id, name, description, picture, price, total_stock, available_stock, reserved_stock, low_stock_threshold, allow_negative, created_at, updated_at, version
	//  FROM products.products
	//  WHERE name ILIKE '%' || $1 || '%'
	SearchProducts(ctx context.Context, dollar_1 *string) ([]ProductsProducts, error)
	// 预留库存（下单时）
	//
	//  UPDATE products.products
	//  SET reserved_stock = reserved_stock + 2
	//  WHERE id = 1
	//  RETURNING id, name, description, picture, price, total_stock, available_stock, reserved_stock, low_stock_threshold, allow_negative, created_at, updated_at, version
	UpdateProductsReservedStock(ctx context.Context) (ProductsProducts, error)
}

var _ Querier = (*Queries)(nil)

