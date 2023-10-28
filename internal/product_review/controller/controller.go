package controller

import (
	"context"
	"errors"
	"fmt"

	core "github.com/vincentcent6/product-catalogue/internal/product_review/core"
)

type (
	Controller interface {
		CreateProductReview(ctx context.Context, input ProductReviewInput) (result ProductReviewInput, err error)
	}
	ctrl struct{}

	ProductReviewInput struct {
		ProductID     int64  `json:"product_id,omitempty"`
		Rating        int    `json:"rating,omitempty"`
		ReviewComment string `json:"review_comment,omitempty"`
	}
)

var (
	prdReview core.ProductReview

	ErrProductIDLTEZero = errors.New("Product ID is less than or equals zero")
	ErrRatingLTEZero    = errors.New("Rating is zero")
)

func New() Controller {
	err := prepare()
	if err != nil {
		fmt.Println("Product New Controller - failed to prepare")
	}
	return &ctrl{}
}

func prepare() (err error) {
	if prdReview == nil {
		prdReview = core.NewProductReview()
	}
	return nil
}

func (p ProductReviewInput) ValidateComplete() error {
	if p.ProductID <= 0 {
		return ErrProductIDLTEZero
	} else if p.Rating <= 0 {
		return ErrRatingLTEZero
	}
	return nil
}
