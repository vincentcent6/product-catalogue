package controller

import (
	"context"
	"errors"
	"fmt"

	"github.com/vincentcent6/product-catalogue/internal/core"
)

type (
	Controller interface {
		CreateProduct(ctx context.Context, input ProductInput) (result ProductInput, err error)
		UpdateProduct(ctx context.Context, input UpdateProductInput) (err error)
		GetProduct(ctx context.Context, input GetProductInput) (result GetProductResult, err error)
		GetProducts(ctx context.Context, input GetProductInput) (result []GetProductResult, err error)
	}
	ctrl struct{}

	ProductInput struct {
		SKU         string         `json:"sku,omitempty"`
		Title       string         `json:"title,omitempty"`
		Description string         `json:"description,omitempty"`
		Category    string         `json:"category,omitempty"`
		Etalase     string         `json:"etalase,omitempty"`
		Images      []ProductImage `json:"images,omitempty"`
		Weight      float64        `json:"weight,omitempty"`
		Price       int64          `json:"price,omitempty"`
	}

	ProductImage struct {
		URL         string `json:"url,omitempty"`
		Description string `json:"description,omitempty"`
	}

	UpdateProductInput struct {
		ProductID int64        `json:"product_id,omitempty"`
		Attribute ProductInput `json:"attribute,omitempty"`
	}

	GetProductInput struct {
		ProductID int64             `json:"product_id,omitempty"`
		Filter    *GetProductFilter `json:"filter,omitempty"`
		Sort      *GetProductSort   `json:"sort,omitempty"`
	}

	GetProductResult struct {
		ProductID int64        `json:"product_id,omitempty"`
		Attribute ProductInput `json:"attribute,omitempty"`
	}

	GetProductFilter struct {
		SKU         string `json:"sku,omitempty"`
		Title       string `json:"title,omitempty"`
		Description string `json:"description,omitempty"`
		Category    string `json:"category,omitempty"`
	}

	GetProductSort struct {
		CreateTime int `json:"sort_create_time,omitempty"`
	}
)

var (
	prd core.Product

	ErrSKUIsEmpty         = errors.New("SKU is empty")
	ErrTitleIsEmpty       = errors.New("Title is empty")
	ErrCategoryIsEmpty    = errors.New("Category is empty")
	ErrDescriptionIsEmpty = errors.New("Description is empty")
	ErrEtalaseIsEmpty     = errors.New("Etalase is empty")
	ErrImagesIsEmpty      = errors.New("Images is empty")
	ErrWeightLessThanZero = errors.New("Weight is less than zero")
	ErrPriceLTEZero       = errors.New("Price is less than or equals zero")
	ErrAnyInvalidValue    = errors.New("An attribute has invalid value")
)

func New() Controller {
	err := prepare()
	if err != nil {
		fmt.Println("Product New Controller - failed to prepare")
	}
	return &ctrl{}
}

func prepare() (err error) {
	if prd == nil {
		prd = core.NewProduct()
	}
	return nil
}

func (p ProductInput) ValidateComplete() error {
	if p.SKU == "" {
		return ErrSKUIsEmpty
	} else if p.Title == "" {
		return ErrTitleIsEmpty
	} else if p.Description == "" {
		return ErrDescriptionIsEmpty
	} else if p.Category == "" {
		return ErrCategoryIsEmpty
	} else if p.Etalase == "" {
		return ErrEtalaseIsEmpty
	} else if len(p.Images) == 0 {
		return ErrImagesIsEmpty
	} else if p.Weight < 0 {
		return ErrWeightLessThanZero
	} else if p.Price <= 0 {
		return ErrPriceLTEZero
	}
	return nil
}
