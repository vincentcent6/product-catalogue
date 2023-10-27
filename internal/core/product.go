package core

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type (
	// Product is an interface for product object
	Product interface {
		CreateProduct(ctx context.Context, data Data) error
	}
	product struct{}

	Data struct {
		SKU         int64          `db:"sku" json:"sku,omitempty"`
		Title       string         `db:"title" json:"title,omitempty"`
		Description string         `db:"description" json:"description,omitempty"`
		Category    string         `db:"category" json:"category,omitempty"`
		Etalase     string         `db:"etalase" json:"etalase,omitempty"`
		Images      []ProductImage `db:"images" json:"images,omitempty"`
		Weight      float64        `db:"weight" json:"weight,omitempty"`
		Price       float64        `db:"price" json:"price,omitempty"`
	}

	ProductImage struct {
		URL         string `db:"url" json:"url,omitempty"`
		Description string `db:"description" json:"description,omitempty"`
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

func (img ProductImage) Scan(val interface{}) error {
	switch v := val.(type) {
	case []byte:
		json.Unmarshal(v, &img)
		return nil
	case string:
		json.Unmarshal([]byte(v), &img)
		return nil
	default:
		return errors.New(fmt.Sprintf("Unsupported type: %T", v))
	}
}

func (img ProductImage) Value() (driver.Value, error) {
	return json.Marshal(img)
}
