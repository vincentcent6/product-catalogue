package controller

import (
	"context"

	core "github.com/vincentcent6/product-catalogue/internal/product_review/core"
)

func (p *ctrl) CreateProductReview(ctx context.Context, input ProductReviewInput) (result ProductReviewInput, err error) {
	if err = input.ValidateComplete(); err != nil {
		return ProductReviewInput{}, err
	}
	if err = prdReview.CreateProductReview(ctx, convertData(input)); err != nil {
		return ProductReviewInput{}, err
	}
	return input, nil
}

func convertData(input ProductReviewInput) core.Data {
	return core.Data{
		ProductID:     input.ProductID,
		Rating:        input.Rating,
		ReviewComment: input.ReviewComment,
	}
}
