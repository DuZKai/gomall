package model

import (
	"context"

	"gorm.io/gorm"
)

type Category struct {
	Base
	Name        string `json:"name"`
	Description string `json:"description"`

	Products []Product `json:"product" gorm:"many2many:product_category"`
}

func (c Category) TableName() string {
	return "category"
}

type CategoryQuery struct {
	ctx context.Context
	db  *gorm.DB
}

func (c CategoryQuery) GetProductsByCategoryName(name string) (categories []Category, err error) {
	err = c.db.WithContext(c.ctx).Model(&Category{}).Where(&Category{Name: name}).Preload("Products").Find(&categories).Error
	return
}

// 用于创建 CategoryQuery 实例
func NewCategoryQuery(ctx context.Context, db *gorm.DB) *CategoryQuery {
	return &CategoryQuery{
		ctx: ctx,
		db:  db,
	}
}

func (c CategoryQuery) GetProductsByCategoryNameAndPage(name string, page, pageSize int32) (categories []Category, err error) {
	offset := (page - 1) * pageSize
	err = c.db.WithContext(c.ctx).
		Model(&Category{}).
		Where(&Category{Name: name}).
		// 加载带分页的 Products
		Preload("Products", func(db *gorm.DB) *gorm.DB {
			return db.Offset(int(offset)).Limit(int(pageSize))
		}).
		Find(&categories).Error
	return
}
