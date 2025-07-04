package model

import (
	"context"
	"gorm.io/gorm"
)

type Product struct {
	Base
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Picture     string  `json:"picture"`
	Price       float32 `json:"price"`

	Categories []Category `json:"categories" gorm:"many2many:product_category"`
}

func (p Product) TableName() string {
	return "product"
}

type ProductQuery struct {
	ctx context.Context
	db  *gorm.DB
}

func (p ProductQuery) GetById(productId int) (product Product, err error) {
	err = p.db.WithContext(p.ctx).Model(&Product{}).First(&product, productId).Error
	return
}

func (p ProductQuery) SearchProducts(q string) (products []*Product, err error) {
	err = p.db.WithContext(p.ctx).Model(&Product{}).Find(&products, "name like ? or description like ?",
		"%"+q+"%", "%"+q+"%",
	).Error
	return
}

// 用于创建 ProductQuery 实例
func NewProductQuery(ctx context.Context, db *gorm.DB) *ProductQuery {
	return &ProductQuery{
		ctx: ctx,
		db:  db,
	}
}

// 根据id修改product
func (p ProductQuery) UpdateProduct(productId int, product *Product) (err error) {
	err = p.db.WithContext(p.ctx).Model(&Product{}).Where("id = ?", productId).Updates(product).Error
	if err != nil {
		return err
	}
	return nil
}

// 查询所有ID
func (p ProductQuery) GetAllId() (ids []*uint32, err error) {
	err = p.db.WithContext(p.ctx).Model(&Product{}).Select("id").Find(&ids).Error
	return
}
