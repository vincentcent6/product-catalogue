package controller

import (
	"context"

	"github.com/vincentcent6/product-catalogue/internal/core"
)

func (p *ctrl) GetProduct(ctx context.Context, input GetProductInput) (result GetProductResult, err error) {
	res, err := prd.GetProduct(ctx, convertGetData(input))
	if err != nil {
		return GetProductResult{}, err
	}

	return parseGetResult(res), err
}

func convertGetData(input GetProductInput) core.GetInput {
	return core.GetInput{
		ProductID: input.ProductID,
	}
}

func parseGetResult(prd core.ProductData) GetProductResult {
	images := []ProductImage{}
	for _, img := range prd.Images {
		images = append(images, ProductImage{
			URL:         img.URL,
			Description: img.Description,
		})
	}

	return GetProductResult{
		ProductID: prd.ProductID,
		Attribute: ProductInput{
			SKU:         prd.SKU,
			Title:       prd.Title,
			Description: prd.Description,
			Category:    prd.Category,
			Etalase:     prd.Etalase,
			Images:      images,
			Weight:      prd.Weight,
			Price:       prd.Price,
		},
	}
}
