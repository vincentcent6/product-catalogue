package core

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type (
	// Product is an interface for product object
	Product interface {
		CreateProduct(ctx context.Context, data Data) error
		UpdateProduct(ctx context.Context, input UpdateInput) error
		GetProduct(ctx context.Context, inp GetInput) (ProductData, error)
	}
	product struct{}

	ProductImages []ProductImage

	Data struct {
		SKU         int64         `db:"sku" json:"sku,omitempty"`
		Title       string        `db:"title" json:"title,omitempty"`
		Description string        `db:"description" json:"description,omitempty"`
		Category    string        `db:"category" json:"category,omitempty"`
		Etalase     string        `db:"etalase" json:"etalase,omitempty"`
		Images      ProductImages `db:"images" json:"images,omitempty"`
		Weight      float64       `db:"weight" json:"weight,omitempty"`
		Price       int64         `db:"price" json:"price,omitempty"`
	}

	UpdateInput struct {
		ProductID int64 `db:"product_id" json:"product_id,omitempty"`
		Data      Data
	}

	ProductData struct {
		ProductID   int64         `db:"product_id" json:"product_id,omitempty"`
		SKU         int64         `db:"sku" json:"sku,omitempty"`
		Title       string        `db:"title" json:"title,omitempty"`
		Description string        `db:"description" json:"description,omitempty"`
		Category    string        `db:"category" json:"category,omitempty"`
		Etalase     string        `db:"etalase" json:"etalase,omitempty"`
		Images      ProductImages `db:"images" json:"images,omitempty"`
		Weight      float64       `db:"weight" json:"weight,omitempty"`
		Price       int64         `db:"price" json:"price,omitempty"`
	}

	ProductImage struct {
		URL         string `db:"url" json:"url,omitempty"`
		Description string `db:"description" json:"description,omitempty"`
	}

	GetInput struct {
		ProductID int64 `db:"product_id" json:"product_id,omitempty"`
	}
)

func (p *product) CreateProduct(ctx context.Context, data Data) error {
	imgJson, err := json.Marshal(data.Images)
	if err != nil {
		return err
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(qCreateProduct, data.SKU, data.Title, data.Description, data.Category, data.Etalase, imgJson, data.Weight, data.Price)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (p *product) UpdateProduct(ctx context.Context, input UpdateInput) error {
	data := input.Data
	params := []string{"", "", "", "", "", "", "", ""}

	imgJson, err := json.Marshal(data.Images)
	if err != nil {
		return err
	}
	// check if the input is set
	if data.SKU != 0 {
		params[0] = fmt.Sprintf("sku=%d", data.SKU)
	}
	if data.Title != "" {
		params[1] = fmt.Sprintf("title='%s'", data.Title)
	}
	if data.Description != "" {
		params[2] = fmt.Sprintf("description='%s'", data.Description)
	}
	if data.Category != "" {
		params[3] = fmt.Sprintf("category='%s'", data.Category)
	}
	if data.Etalase != "" {
		params[4] = fmt.Sprintf("etalase='%s'", data.Etalase)
	}
	if len(data.Images) > 0 {
		params[5] = fmt.Sprintf("images='%v'", string(imgJson))
	}
	if data.Weight != 0 {
		params[6] = fmt.Sprintf("weight=%f", data.Weight)
	}
	if data.Price != 0 {
		params[7] = fmt.Sprintf("price=%d", data.Price)
	}
	paramsStr := strings.Join(params, ", ")

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	res, err := tx.Exec(fmt.Sprintf(qUpdateProduct, paramsStr), input.ProductID)
	if err != nil {
		return err
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return ErrDataNotFound
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

// for now, get product will return 1 result
func (p *product) GetProduct(ctx context.Context, inp GetInput) (ProductData, error) {
	params := []string{""}
	if inp.ProductID != 0 {
		params[0] = fmt.Sprintf("product_id=%d", inp.ProductID)
	}
	paramsStr := strings.Join(params, ", ")

	prd := []ProductData{}
	err := db.Select(&prd, fmt.Sprintf(qGetProduct, paramsStr))
	if err != nil {
		fmt.Println(err)
		return ProductData{}, err
	}
	if len(prd) == 0 {
		return ProductData{}, ErrDataNotFound
	}
	return prd[0], nil
}

func (images ProductImages) Value() (driver.Value, error) {
	return json.Marshal(images)
}

func (images *ProductImages) Scan(val interface{}) (err error) {
	switch v := val.(type) {
	case []byte:
		return json.Unmarshal(v, images)
	case string:
		return json.Unmarshal([]byte(v), &images)
	default:
		return errors.New(fmt.Sprintf("Unsupported type: %T", v))
	}
}
