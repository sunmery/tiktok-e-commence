// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package models

import (
	"context"
)

type Querier interface {
	//CreateCategory
	//
	//  /*
	//  参数说明：
	//  - parent_id: 父分类ID（0表示根）
	//  - name: 分类名称
	//  - sort_order: 排序序号
	//
	//  操作流程：
	//  1. 检查父分类是否存在
	//  2. 插入新分类并生成ltree路径
	//  3. 维护闭包表关系
	//  */
	//
	//  WITH root_check AS (
	//    INSERT INTO categories.categories (id, parent_id, level, path, name, sort_order, is_leaf)
	//    VALUES (0, 0, 1, 'root'::public.ltree, 'Root', 0, FALSE)
	//    ON CONFLICT (id) DO NOTHING
	//  ),
	//  parent_info AS (
	//    SELECT
	//      COALESCE(c.id, 0) AS effective_parent_id,
	//      COALESCE(c.path, 'root'::public.ltree) AS parent_path,
	//      COALESCE(c.level, 0) AS parent_level
	//    FROM (SELECT $1::BIGINT AS pid) AS input
	//    LEFT JOIN categories.categories c ON c.id = input.pid
	//  ),
	//  level_validation AS (
	//    SELECT
	//      effective_parent_id,
	//      parent_path,
	//      CASE
	//        WHEN parent_level >= 4 THEN NULL -- 父节点已经是4层，不允许新增子节点
	//        ELSE parent_level + 1
	//      END AS new_level
	//    FROM parent_info
	//  ),
	//  insert_main AS (
	//    INSERT INTO categories.categories (
	//      parent_id, level, path, name, sort_order, is_leaf
	//    ) SELECT
	//      lv.effective_parent_id,
	//      lv.new_level,
	//      CASE
	//        WHEN lv.parent_path OPERATOR(public.=) 'root'::public.ltree
	//        THEN lv.parent_path || ('node_' || gen_random_uuid())::public.ltree
	//        ELSE lv.parent_path || (REPLACE(gen_random_uuid()::text, '-', '_'))::public.ltree
	//      END,
	//      $2,   -- Name 参数
	//      $3,   -- SortOrder 参数
	//      CASE WHEN lv.new_level = 4 THEN TRUE ELSE FALSE END -- 第四层为叶子节点
	//    FROM level_validation lv
	//    WHERE lv.new_level IS NOT NULL
	//    RETURNING id, parent_id, level, path, name, sort_order, is_leaf, created_at, updated_at
	//  ),
	//  update_parent_leaf AS (
	//    UPDATE categories.categories
	//    SET is_leaf = FALSE
	//    WHERE id = (SELECT effective_parent_id FROM parent_info)
	//    AND is_leaf = TRUE
	//  )
	//  INSERT INTO categories.category_closure (ancestor, descendant, depth)
	//  SELECT
	//    cc.ancestor,
	//    im.id,
	//    cc.depth + 1
	//  FROM insert_main im
	//  JOIN categories.category_closure cc ON cc.descendant = im.parent_id
	//  UNION ALL
	//  SELECT
	//    im.id,
	//    im.id,
	//    0
	//  FROM insert_main im
	//  RETURNING descendant
	CreateCategory(ctx context.Context, arg CreateCategoryParams) (int64, error)
	//DeleteCategory
	//
	//  WITH deleted AS (
	//      DELETE FROM categories.categories
	//          WHERE id = $1
	//          RETURNING path
	//  )
	//  DELETE FROM categories.category_closure
	//  WHERE descendant IN (
	//      SELECT descendant
	//      FROM categories.category_closure
	//      WHERE ancestor = $1
	//  )
	DeleteCategory(ctx context.Context, id *int64) error
	// 确保深度不超过 3
	// 删除指定分类及其所有后代节点的闭包关系
	//
	//
	//  DELETE FROM categories.category_closure
	//  WHERE descendant IN (
	//      SELECT descendant
	//      FROM categories.category_closure
	//      WHERE ancestor = $1
	//  )
	DeleteClosureRelations(ctx context.Context, categoryID *int64) error
	//GetCategoryByID
	//
	//  SELECT id, parent_id, level, path, name, sort_order, is_leaf, created_at, updated_at FROM categories.categories
	//  WHERE id = $1 LIMIT 1
	GetCategoryByID(ctx context.Context, id int64) (CategoriesCategories, error)
	//GetCategoryPath
	//
	//  SELECT ancestor.id, ancestor.parent_id, ancestor.level, ancestor.path, ancestor.name, ancestor.sort_order, ancestor.is_leaf, ancestor.created_at, ancestor.updated_at
	//  FROM categories.category_closure cc
	//           JOIN categories.categories ancestor ON cc.ancestor = ancestor.id
	//  WHERE cc.descendant = $1
	//  ORDER BY cc.depth DESC
	GetCategoryPath(ctx context.Context, categoryID int64) ([]CategoriesCategories, error)
	//GetClosureRelations
	//
	//  SELECT ancestor, descendant, depth FROM categories.category_closure
	//  WHERE descendant = $1
	GetClosureRelations(ctx context.Context, categoryID int64) ([]CategoriesCategoryClosure, error)
	//GetLeafCategories
	//
	//  SELECT id, parent_id, level, path, name, sort_order, is_leaf, created_at, updated_at FROM categories.categories
	//  WHERE is_leaf = TRUE AND level = 4
	GetLeafCategories(ctx context.Context) ([]CategoriesCategories, error)
	//GetSubTree
	//
	//  /*
	//  级联删除策略：
	//  1. 根据闭包表找到所有后代节点
	//  2. 删除所有相关闭包关系
	//  */
	//
	//
	//  SELECT c.id, c.parent_id, c.level, c.path, c.name, c.sort_order, c.is_leaf, c.created_at, c.updated_at
	//  FROM categories.categories c
	//  WHERE c.path <@ (SELECT path FROM categories.categories WHERE id = $1)
	//  ORDER BY c.path
	GetSubTree(ctx context.Context, rootID *int64) ([]CategoriesCategories, error)
	//UpdateCategoryName
	//
	//  UPDATE categories.categories
	//  SET name = $1, updated_at = NOW()
	//  WHERE id = $2
	UpdateCategoryName(ctx context.Context, arg UpdateCategoryNameParams) error
	//UpdateClosureDepth
	//
	//  UPDATE categories.category_closure
	//  SET depth = depth + $1
	//  WHERE descendant IN (
	//      SELECT descendant
	//      FROM categories.category_closure
	//      WHERE ancestor = $2
	//  )
	//  AND depth + $1 <= 3
	UpdateClosureDepth(ctx context.Context, arg UpdateClosureDepthParams) error
	// 更新父分类的叶子节点状态
	//
	//  UPDATE categories.categories
	//  SET
	//      is_leaf = NOT EXISTS (
	//          SELECT 1
	//          FROM categories
	//          WHERE parent_id = $1
	//          LIMIT 1
	//      ),
	//      updated_at = NOW()
	//  WHERE id = $1
	UpdateParentLeafStatus(ctx context.Context, parentID *int64) error
}

var _ Querier = (*Queries)(nil)
